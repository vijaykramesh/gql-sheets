package model

import (
	"errors"
	"fmt"
	"github.com/xuri/efp"
	"gorm.io/gorm"
	"regexp"
	"strconv"
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
