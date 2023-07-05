package model

import (
	"errors"
	"fmt"
	"github.com/xuri/efp"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"strings"
)

type Cell struct {
	gorm.Model
	SpreadsheetID string       `json:"spreadsheetId"`
	Spreadsheet   *Spreadsheet `json:"spreadsheet"`
	RawValue      string       `json:"rawValue"`
	ComputedValue string       `json:"computedValue,omitempty"`
	RowIndex      int          `json:"rowIndex"`
	ColumnIndex   int          `json:"columnIndex"`
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

		// =SUM(A1:A3) style formula, compute the value
		if len(tokens) == 3 && tokens[0].TType == "Function" && tokens[0].TSubType == "Start" && tokens[0].TValue == "SUM" && tokens[1].TType == "Operand" && tokens[1].TSubType == "Range" && tokens[2].TType == "Function" && tokens[2].TSubType == "Stop" && tokens[2].TValue == "" {
			fmt.Println("SUM formula")
			fmt.Println(tokens)
			fmt.Println(tokens[1].TValue)
			// split tokens[1].TValue on : to get start and end cell
			startCell := strings.Split(tokens[1].TValue, ":")[0]
			endCell := strings.Split(tokens[1].TValue, ":")[1]

			startColumnIndex, startRowIndex, err := columnAndRowIndexFromCode(startCell)
			if err != nil {
				return "", err
			}

			endColumnIndex, endRowIndex, err := columnAndRowIndexFromCode(endCell)
			if err != nil {
				return "", err
			}

			var sum int64
			for _, cell := range otherCells {
				if cell.ColumnIndex >= startColumnIndex && cell.ColumnIndex <= endColumnIndex && cell.RowIndex >= startRowIndex && cell.RowIndex <= endRowIndex {
					cellValue, err := strconv.ParseInt(cell.ComputedValue, 10, 64)
					if err != nil {
						return "", err
					}
					sum += cellValue
				}
			}
			return strconv.FormatInt(sum, 10), nil
		}
		return c.RawValue, nil
	}
	return c.RawValue, nil
}
func (c *Cell) FindDependentCells(otherCells []Cell) ([]Cell, error) {
	var dependentCells []Cell
	for _, cell := range otherCells {
		tokens, err := cell.parseRawValue()

		cToken := columnCodeFromColumnIndex(c.ColumnIndex) + strconv.Itoa(c.RowIndex+1)
		if err == nil && len(tokens) == 1 && tokens[0].TType == "Operand" && tokens[0].TSubType == "Range" && tokens[0].TValue == cToken {
			dependentCells = append(dependentCells, cell)
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
