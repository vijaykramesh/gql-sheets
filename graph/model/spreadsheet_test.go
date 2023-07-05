package model

import "testing"

func TestValidateRowAndColumnIndexes(t *testing.T) {
	// Create a test spreadsheet
	spreadsheet := Spreadsheet{
		RowCount:    5,
		ColumnCount: 5,
	}

	// Test case 1: Valid indexes
	err := ValidateRowAndColumnIndexes(spreadsheet, 3, 2)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Test case 2: Row index greater than row count
	err = ValidateRowAndColumnIndexes(spreadsheet, 6, 2)
	expectedError := "row index 6 is greater than row count 5"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got: '%v'", expectedError, err)
	}

	// Test case 3: Column index greater than column count
	err = ValidateRowAndColumnIndexes(spreadsheet, 3, 6)
	expectedError = "column index 6 is greater than column count 5"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got: '%v'", expectedError, err)
	}
}
