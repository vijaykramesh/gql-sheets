package model

import (
	"errors"
	"fmt"
	"github.com/WinterYukky/gorm-extra-clause-plugin/exclause"
	"github.com/vijaykramesh/gql-sheets/graph/common"
	"github.com/xuri/efp"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Cell struct {
	gorm.Model
	SpreadsheetID string       `json:"spreadsheetId"`
	Spreadsheet   *Spreadsheet `json:"spreadsheet"`
	RawValue      string       `json:"rawValue"`
	ComputedValue string       `json:"computedValue,omitempty"`
	RowIndex      int          `json:"rowIndex"`
	ColumnIndex   int          `json:"columnIndex"`
	Version       uint64       `json:"version"`
}

func (c *Cell) parseRawValue() ([]efp.Token, error) {
	if len(c.RawValue) > 0 && c.RawValue[0] != '=' {
		return nil, errors.New("raw value does not start with =")
	}
	ps := efp.ExcelParser()
	ps.Parse(c.RawValue)
	return ps.Tokens.Items, nil
}

func (c *Cell) ComputeValueFromRaw(otherCells []Cell) (string, error) {
	// if c.RawValue starts with =, then it is a formula or reference
	if len(c.RawValue) > 0 && c.RawValue[0] == '=' {
		tokens, err := c.parseRawValue()
		if err != nil {
			return "", err
		}

		// =A4 style reference, lookup the computed value of the referenced cell
		if len(tokens) == 1 && tokens[0].TType == "Operand" && tokens[0].TSubType == "Range" {
			computedValue, lookupErr := referenceLookup(tokens, otherCells)
			if lookupErr != nil {
				return "ERROR", lookupErr
			}
			return computedValue, nil
		}
		// =SUM(A1:A3) style formula, compute the value
		if len(tokens) == 3 && tokens[0].TType == "Function" && tokens[0].TSubType == "Start" && tokens[0].TValue == "SUM" && tokens[1].TType == "Operand" && tokens[1].TSubType == "Range" && tokens[2].TType == "Function" && tokens[2].TSubType == "Stop" && tokens[2].TValue == "" {
			computedValue, lookupErr := sumRange(tokens, otherCells)
			if lookupErr != nil {
				return "ERROR", lookupErr
			}
			return computedValue, nil
		}

		// =AVERAGE(A1:A3) style formula, compute the value
		if len(tokens) == 3 && tokens[0].TType == "Function" && tokens[0].TSubType == "Start" && tokens[0].TValue == "AVERAGE" && tokens[1].TType == "Operand" && tokens[1].TSubType == "Range" && tokens[2].TType == "Function" && tokens[2].TSubType == "Stop" && tokens[2].TValue == "" {
			computedValue, lookupErr := averageRange(tokens, otherCells)
			if lookupErr != nil {
				return "ERROR", lookupErr
			}
			return computedValue, nil
		}
		return c.RawValue, nil
	}
	return c.RawValue, nil
}

func (c *Cell) UpdateCellAndDependentCells(context *common.CustomContext, input UpdateCell) (*Cell, error) {
	version := uint64(time.Now().UnixMilli())
	c.RawValue = input.RawValue

	var otherCells []Cell
	err := context.Database.Clauses(exclause.NewWith("cte", context.Database.Table("cells").Select("column_index,row_index,max(version) as version").Group("column_index,row_index"))).Where("spreadsheet_id = ? AND (column_index != ? OR row_index != ?) AND version = (SELECT version FROM cte WHERE column_index = cells.column_index AND row_index = cells.row_index)", c.SpreadsheetID, c.ColumnIndex, c.RowIndex).Find(&otherCells).Error
	if err != nil {
		return nil, fmt.Errorf("error getting cells: %v", err)
	}
	latestVersionSeenForColumnAndRowIndex := make(map[string]uint64)
	for _, otherCell := range otherCells {
		key := fmt.Sprintf("%d-%d", otherCell.ColumnIndex, otherCell.RowIndex)
		if latestVersionSeenForColumnAndRowIndex[key] < otherCell.Version {
			latestVersionSeenForColumnAndRowIndex[key] = otherCell.Version
		}
	}

	onlyLatestVersionOtherCells := make([]Cell, 0, len(otherCells))
	for _, otherCell := range otherCells {
		key := fmt.Sprintf("%d-%d", otherCell.ColumnIndex, otherCell.RowIndex)
		if latestVersionSeenForColumnAndRowIndex[key] == otherCell.Version {
			onlyLatestVersionOtherCells = append(onlyLatestVersionOtherCells, otherCell)
		}
	}
	c.ComputedValue, err = c.ComputeValueFromRaw(onlyLatestVersionOtherCells)
	c.Version = version
	err = context.Database.Omit("id").Create(&c).Error
	// select cell we just saved and set to c
	newC := Cell{}
	err = context.Database.Clauses(exclause.NewWith("cte", context.Database.Table("cells").Select("column_index,row_index,max(version) as version").Group("column_index,row_index"))).Where("spreadsheet_id = ? AND column_index = ? and row_index = ? AND version = (SELECT version FROM cte WHERE column_index = cells.column_index AND row_index = cells.row_index)", c.SpreadsheetID, c.ColumnIndex, c.RowIndex).First(&newC).Error
	if err != nil {
		return nil, fmt.Errorf("error updating cell: %v", err)
	}

	onlyLatestVersionOtherCells = append(onlyLatestVersionOtherCells, newC)
	dependentCells, err := c.FindDependentCells(onlyLatestVersionOtherCells)
	for _, dependentCell := range dependentCells {
		dependentCell.ComputedValue, err = dependentCell.ComputeValueFromRaw(onlyLatestVersionOtherCells)
		if err != nil {
			return nil, err
		}
		dependentCell.Version = version
		err = context.Database.Omit("id").Create(&dependentCell).Error
		if err != nil {
			return nil, fmt.Errorf("error updating cell: %v", err)
		}

		// todo change this to recurse and support infinite reference depth
		filteredOtherCells := []Cell{}
		for _, otherCell := range onlyLatestVersionOtherCells {
			if otherCell.ID != dependentCell.ID {
				filteredOtherCells = append(filteredOtherCells, otherCell)
			}
		}
		filteredOtherCells = append(filteredOtherCells, dependentCell)
		dependentCells, err = dependentCell.FindDependentCells(filteredOtherCells)
		for _, dc := range dependentCells {
			dc.ComputedValue, err = dc.ComputeValueFromRaw(filteredOtherCells)
			if err != nil {
				return nil, err
			}
			dc.Version = version
			err = context.Database.Omit("id").Create(&dc).Error
		}

	}

	return c, nil
}

func referenceLookup(tokens []efp.Token, otherCells []Cell) (string, error) {
	lookupCell := tokens[0].TValue
	columnIndex, rowIndex, err := columnAndRowIndexFromCode(lookupCell)
	if err != nil {
		return "", err
	}

	var matchedCell *Cell
	for _, cell := range otherCells {
		if cell.ColumnIndex == columnIndex && cell.RowIndex == rowIndex {
			matchedCell = &cell
			break
		}
	}
	if matchedCell == nil {
		return "", fmt.Errorf("could not find cell %s", lookupCell)
	}
	return matchedCell.ComputedValue, nil
}
func averageRange(tokens []efp.Token, otherCells []Cell) (string, error) {
	// Compute average
	var sum float64
	var count int64
	for _, cell := range otherCells {
		check, err := checkIfCellInRange(&cell, tokens[1].TValue)
		if err != nil {
			return "", err
		}
		if check {
			cellValue, err := strconv.ParseFloat(cell.ComputedValue, 64)
			if err == nil {
				sum += cellValue
				count++
			}
		}
	}
	average := sum / float64(count)
	return strconv.FormatFloat(average, 'f', -1, 64), nil
}

func sumRange(tokens []efp.Token, otherCells []Cell) (string, error) {

	var sum int64
	for _, cell := range otherCells {
		check, e := checkIfCellInRange(&cell, tokens[1].TValue)
		if e != nil {
			return "", e
		}
		if check {
			cellValue, err := strconv.ParseInt(cell.ComputedValue, 10, 64)
			if err == nil {
				sum += cellValue
			}
		}
	}
	return strconv.FormatInt(sum, 10), nil
}

func checkIfCellInRange(c *Cell, tvalue string) (bool, error) {
	if len(strings.Split(tvalue, ":")) != 2 {
		return false, fmt.Errorf("invalid range %s", tvalue)
	}
	startCell := strings.Split(tvalue, ":")[0]
	endCell := strings.Split(tvalue, ":")[1]

	startColumnIndex, startRowIndex, err := columnAndRowIndexFromCode(startCell)
	if err != nil {
		return false, err
	}
	endColumnIndex, endRowIndex, err := columnAndRowIndexFromCode(endCell)
	if err != nil {
		return false, err
	}
	if c.ColumnIndex >= startColumnIndex && c.ColumnIndex <= endColumnIndex && c.RowIndex >= startRowIndex && c.RowIndex <= endRowIndex {
		return true, nil
	}
	return false, nil
}

func (c *Cell) FindDependentCells(otherCells []Cell) ([]Cell, error) {
	var dependentCells []Cell
	for _, cell := range otherCells {
		tokens, err := cell.parseRawValue()

		// =A4 style reference, check if c is the referenced cell and if so add cell to dependent cells
		cToken := columnCodeFromColumnIndex(c.ColumnIndex) + strconv.Itoa(c.RowIndex+1)
		if err == nil && len(tokens) == 1 && tokens[0].TType == "Operand" && tokens[0].TSubType == "Range" && tokens[0].TValue == cToken {
			dependentCells = append(dependentCells, cell)
		}
		// =SUM(A1:A3) style formula, check if c is in the range and if so add cell to dependent cells
		if err == nil && len(tokens) == 3 && tokens[0].TType == "Function" && tokens[0].TSubType == "Start" && (tokens[0].TValue == "SUM" || tokens[0].TValue == "AVERAGE") && tokens[1].TType == "Operand" && tokens[1].TSubType == "Range" && tokens[2].TType == "Function" && tokens[2].TSubType == "Stop" && tokens[2].TValue == "" {
			check, _ := checkIfCellInRange(c, tokens[1].TValue)
			if check {
				dependentCells = append(dependentCells, cell)
			}
		}
	}
	return dependentCells, nil

}

func columnAndRowIndexFromCode(code string) (int, int, error) {
	// Use regular expression to split the code into letters and numbers
	re := regexp.MustCompile(`^([A-Z]+)(\d+)$`)
	matches := re.FindStringSubmatch(code)
	if len(matches) != 3 {
		return -1, -1, fmt.Errorf("invalid code format: %s", code)
	}

	letters := matches[1]
	numbers := matches[2]

	// Compute the column index based on letters
	column := 0
	for i := 0; i < len(letters); i++ {
		char := letters[i]
		if char < 'A' || char > 'Z' {
			return -1, -1, fmt.Errorf("invalid column code: %s", code)
		}
		column = (i)*26 + int(char-'A'+1)
	}

	// Compute the row index based on numbers
	row, err := strconv.Atoi(numbers)
	if err != nil {
		return -1, -1, fmt.Errorf("invalid row index: %s", numbers)
	}

	return column - 1, row - 1, nil
}

func columnCodeFromColumnIndex(columnIndex int) string {
	if columnIndex < 0 {
		return ""
	}

	column := ""
	for columnIndex >= 0 {
		remainder := columnIndex % 26
		column = string(rune(remainder+'A')) + column
		columnIndex = (columnIndex / 26) - 1

		if columnIndex < 0 {
			break
		}
	}

	return column
}
