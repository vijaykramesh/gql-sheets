package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/xuri/efp"
	"testing"
)

func TestParseRawValue(t *testing.T) {
	t.Run("reference", func(t *testing.T) {
		cell := Cell{
			RawValue: "=A2",
		}

		tokens, err := cell.parseRawValue()

		assert.NoError(t, err)
		assert.NotNil(t, tokens)
		assert.Equal(t, tokens[0].TType, "Operand")
		assert.Equal(t, tokens[0].TValue, "A2")
		assert.Equal(t, tokens[0].TSubType, "Range")

	})
	t.Run("simple formula", func(t *testing.T) {
		cell := Cell{
			RawValue: "=SUM(A1:A3)",
		}

		tokens, err := cell.parseRawValue()

		assert.NoError(t, err)
		assert.NotNil(t, tokens)
		assert.Equal(t, 3, len(tokens))
		// Additional assertions for token values or types can be added as needed
	})
	t.Run("Nested Formula", func(t *testing.T) {
		cell := Cell{
			RawValue: "=SUM(A1:A3) * 2 + MAX(B1:B5)",
		}

		tokens, err := cell.parseRawValue()

		assert.NoError(t, err)
		assert.NotNil(t, tokens)
		assert.Equal(t, 9, len(tokens))
		// Additional assertions for token values or types can be added as needed
	})

	t.Run("Function with Arguments", func(t *testing.T) {
		cell := Cell{
			RawValue: "=IF(A1 > B1, \"Greater\", \"Less or Equal\")",
		}

		tokens, err := cell.parseRawValue()

		assert.NoError(t, err)
		assert.NotNil(t, tokens)
		assert.Equal(t, 9, len(tokens))
		// Additional assertions for token values or types can be added as needed
	})

	t.Run("Garbage input", func(t *testing.T) {
		cell := Cell{
			RawValue: "=asdfadsfsaeeeee",
		}

		tokens, err := cell.parseRawValue()

		assert.NoError(t, err)
		assert.NotNil(t, tokens)
		assert.Equal(t, 1, len(tokens))
		// Additional assertions for token values or types can be added as needed
	})

	t.Run("no Function", func(t *testing.T) {
		cell := Cell{
			RawValue: "42",
		}

		tokens, err := cell.parseRawValue()

		assert.Errorf(t, err, "raw value does not start with =")
		assert.Nil(t, tokens)
		// Additional assertions for token values or types can be added as needed
	})

}

func TestColumnAndRowIndexFromCode(t *testing.T) {
	column, row, _ := columnAndRowIndexFromCode("A1")
	assert.Equal(t, 0, column)
	assert.Equal(t, 0, row)

	column, row, _ = columnAndRowIndexFromCode("C7")
	assert.Equal(t, 2, column)
	assert.Equal(t, 6, row)

	column, row, _ = columnAndRowIndexFromCode("B10")
	assert.Equal(t, 1, column)
	assert.Equal(t, 9, row)

	column, row, _ = columnAndRowIndexFromCode("Z99")
	assert.Equal(t, 25, column)
	assert.Equal(t, 98, row)

	column, row, _ = columnAndRowIndexFromCode("AB77")
	assert.Equal(t, 27, column)
	assert.Equal(t, 76, row)

	column, row, _ = columnAndRowIndexFromCode("ZZ100")
	assert.Equal(t, 51, column)
	assert.Equal(t, 99, row)

	// Invalid column code: empty string
	column, row, err := columnAndRowIndexFromCode("")
	assert.Error(t, err)
	assert.Equal(t, -1, column)
	assert.Equal(t, -1, row)

	// Invalid column code: single character
	column, row, err = columnAndRowIndexFromCode("A")
	assert.Error(t, err)
	assert.Equal(t, -1, column)
	assert.Equal(t, -1, row)

	// Invalid column code: non-uppercase letter
	column, row, err = columnAndRowIndexFromCode("aB77")
	assert.Error(t, err)
	assert.Equal(t, -1, column)
	assert.Equal(t, -1, row)

	// Invalid row index: not a valid integer
	column, row, err = columnAndRowIndexFromCode("ABxyz")
	assert.Error(t, err)
	assert.Equal(t, -1, column)
	assert.Equal(t, -1, row)
}
func TestColumnCodeFromColumnIndex(t *testing.T) {
	// Valid cases
	assert.Equal(t, "A", columnCodeFromColumnIndex(0))
	assert.Equal(t, "B", columnCodeFromColumnIndex(1))
	assert.Equal(t, "Z", columnCodeFromColumnIndex(25))
	assert.Equal(t, "AA", columnCodeFromColumnIndex(26))
	assert.Equal(t, "AB", columnCodeFromColumnIndex(27))
	assert.Equal(t, "AZ", columnCodeFromColumnIndex(51))
	assert.Equal(t, "BA", columnCodeFromColumnIndex(52))

	assert.Equal(t, "AAA", columnCodeFromColumnIndex(702))
	assert.Equal(t, "ABD", columnCodeFromColumnIndex(731))

	// Negative number (invalid case)
	assert.Equal(t, "", columnCodeFromColumnIndex(-1))

	// Error conditions
	assert.Equal(t, "", columnCodeFromColumnIndex(-10))
	assert.Equal(t, "", columnCodeFromColumnIndex(-100))
}

func TestComputeValueFromRaw(t *testing.T) {
	// Prepare test data
	c := &Cell{}
	otherCells := []Cell{
		{ColumnIndex: 0, RowIndex: 0, ComputedValue: "5"},
		{ColumnIndex: 0, RowIndex: 1, ComputedValue: "10"},
		{ColumnIndex: 1, RowIndex: 0, ComputedValue: "words"},
	}

	t.Run("Formula or Reference - Cell Found", func(t *testing.T) {
		// Set up the test case
		c.RawValue = "=A2"

		// Call the function
		result, err := c.ComputeValueFromRaw(otherCells)

		// Check the result
		if err != nil {
			t.Errorf("Expected no error, but got: %s", err.Error())
		}
		expectedValue := "10"
		if result != expectedValue {
			t.Errorf("Expected value: %s, but got: %s", expectedValue, result)
		}
	})

	t.Run("Formula or Reference - Cell Not Found", func(t *testing.T) {
		// Set up the test case
		c.RawValue = "=A3"

		// Call the function
		result, err := c.ComputeValueFromRaw(otherCells)

		// Check the error
		expectedError := "could not find cell A3"
		if err == nil || err.Error() != expectedError {
			t.Errorf("Expected error: %s, but got: %s", expectedError, err)
		}

		// Check the result (should be nil)
		if result != "ERROR" {
			t.Errorf("Expected ERROR result, but got: %s", result)
		}
	})

	t.Run("Non-Formula or Reference", func(t *testing.T) {
		// Set up the test case
		c.RawValue = "123"

		// Call the function
		result, err := c.ComputeValueFromRaw(otherCells)

		// Check the result
		if err != nil {
			t.Errorf("Expected no error, but got: %s", err.Error())
		}
		if result != c.RawValue {
			t.Errorf("Expected value: %s, but got: %s", c.RawValue, result)
		}
	})

	t.Run("SUM Formula - Valid Range", func(t *testing.T) {
		// Set up the test case
		c.RawValue = "=SUM(A1:A2)"

		// Call the function
		result, err := c.ComputeValueFromRaw(otherCells)

		// Check the result
		if err != nil {
			t.Errorf("Expected no error, but got: %s", err.Error())
		}
		expectedValue := "15"
		if result != expectedValue {
			t.Errorf("Expected value: %s, but got: %s", expectedValue, result)
		}
	})
	t.Run("AVERAGE Formula - Valid Range", func(t *testing.T) {
		// Set up the test case
		c.RawValue = "=AVERAGE(A1:A2)"

		// Call the function
		result, err := c.ComputeValueFromRaw(otherCells)

		// Check the result
		if err != nil {
			t.Errorf("Expected no error, but got: %s", err.Error())
		}
		expectedValue := "7.5"
		if result != expectedValue {
			t.Errorf("Expected value: %s, but got: %s", expectedValue, result)
		}
	})
}

func TestReferenceLookup(t *testing.T) {
	tokens := []efp.Token{
		{TValue: "A1"},
	}

	otherCells := []Cell{
		{ColumnIndex: 0, RowIndex: 0, ComputedValue: "Value1"},
		{ColumnIndex: 1, RowIndex: 1, ComputedValue: "Value2"},
	}

	result, err := referenceLookup(tokens, otherCells)
	if err != nil {
		t.Errorf("Error occurred: %s", err.Error())
	}

	expected := "Value1"
	if result != expected {
		t.Errorf("Expected result: %s, but got: %s", expected, result)
	}
}
func TestSumRange(t *testing.T) {
	tokens := []efp.Token{
		{TType: "Function", TSubType: "Start", TValue: "SUM"},
		{TType: "Operand", TSubType: "Range", TValue: "A1:A3"},
		{TType: "Function", TSubType: "Stop", TValue: ""},
	}

	otherCells := []Cell{
		{ColumnIndex: 0, RowIndex: 0, ComputedValue: "10"},
		{ColumnIndex: 0, RowIndex: 1, ComputedValue: "20"},
		{ColumnIndex: 0, RowIndex: 2, ComputedValue: "30"},
		{ColumnIndex: 1, RowIndex: 3, ComputedValue: "40"},
	}

	result, err := sumRange(tokens, otherCells)
	if err != nil {
		t.Errorf("Error occurred: %s", err.Error())
	}

	expected := "60"
	if result != expected {
		t.Errorf("Expected result: %s, but got: %s", expected, result)
	}
}
func TestAverageRange(t *testing.T) {
	tokens := []efp.Token{
		{TType: "Function", TSubType: "Start", TValue: "AVERAGE"},
		{TType: "Operand", TSubType: "Range", TValue: "A1:A3"},
		{TType: "Function", TSubType: "Stop", TValue: ""},
	}

	otherCells := []Cell{
		{ColumnIndex: 0, RowIndex: 0, ComputedValue: "10"},
		{ColumnIndex: 0, RowIndex: 1, ComputedValue: "20"},
		{ColumnIndex: 0, RowIndex: 2, ComputedValue: "30"},
		{ColumnIndex: 1, RowIndex: 1, ComputedValue: "40"},
	}

	result, err := averageRange(tokens, otherCells)
	if err != nil {
		t.Errorf("Error occurred: %s", err.Error())
	}

	expected := "20"
	if result != expected {
		t.Errorf("Expected result: %s, but got: %s", expected, result)
	}
}

func TestCheckIfCellInRange(t *testing.T) {
	c := &Cell{
		ColumnIndex: 1,
		RowIndex:    1,
	}

	tvalue := "A1:B3"

	result, err := checkIfCellInRange(c, tvalue)
	if err != nil {
		t.Errorf("Error occurred: %s", err.Error())
	}

	expected := true
	if result != expected {
		t.Errorf("Expected result: %t, but got: %t", expected, result)
	}

	// Testing case where cell is not within the range
	c = &Cell{
		ColumnIndex: 3,
		RowIndex:    3,
	}

	result, err = checkIfCellInRange(c, tvalue)
	if err != nil {
		t.Errorf("Error occurred: %s", err.Error())
	}

	expected = false
	if result != expected {
		t.Errorf("Expected result: %t, but got: %t", expected, result)
	}

	// Testing case where tvalue is invalid and causes an error
	tvalue = "InvalidRange"

	_, err = checkIfCellInRange(c, tvalue)
	if err == nil {
		t.Error("Expected an error, but got no error")
	}
}

func TestFindDependentCells(t *testing.T) {
	c := Cell{
		ColumnIndex: 0,
		RowIndex:    0,
		RawValue:    "TEST VALUE",
	}

	otherCells := []Cell{
		{ColumnIndex: 1, RowIndex: 0, RawValue: "=A1"},
		{ColumnIndex: 0, RowIndex: 1, RawValue: "=B3"},
		{ColumnIndex: 2, RowIndex: 1, RawValue: "=AVERAGE(A1:A3)"},
		{ColumnIndex: 3, RowIndex: 1, RawValue: "=SUM(A1:C4)"},
	}

	t.Run("Dependent Cells Found", func(t *testing.T) {
		dependentCells, err := c.FindDependentCells(otherCells)
		assert.NoError(t, err)
		assert.NotNil(t, dependentCells)
		assert.Equal(t, 3, len(dependentCells))
		assert.Equal(t, 1, dependentCells[0].ColumnIndex)
		assert.Equal(t, 0, dependentCells[0].RowIndex)
	})

	t.Run("No Dependent Cells", func(t *testing.T) {
		c := Cell{
			ColumnIndex: 4,
			RowIndex:    4,
			RawValue:    "99",
		}

		dependentCells, err := c.FindDependentCells(otherCells)

		assert.NoError(t, err)
		assert.Equal(t, 0, len(dependentCells))
	})

}
