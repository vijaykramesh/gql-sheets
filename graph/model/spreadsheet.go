package model

import (
	"fmt"
	"gorm.io/gorm"
)

type Spreadsheet struct {
	gorm.Model
	Name        string `json:"name"`
	RowCount    int    `json:"rowCount"`
	ColumnCount int    `json:"columnCount"`
}

func ValidateRowAndColumnIndexes(spreadsheet Spreadsheet, rowIndex int, columnIndex int) error {

	if rowIndex >= spreadsheet.RowCount {
		return fmt.Errorf("row index %d is greater than row count %d", rowIndex, spreadsheet.RowCount)
	}
	if columnIndex >= spreadsheet.ColumnCount {
		return fmt.Errorf("column index %d is greater than column count %d", columnIndex, spreadsheet.ColumnCount)
	}
	return nil
}
