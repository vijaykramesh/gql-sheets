package resolvers

import (
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vijaykramesh/gql-sheets/graph/common"
	"github.com/vijaykramesh/gql-sheets/graph/generated"
	"github.com/vijaykramesh/gql-sheets/graph/model"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestCellResolver_ID(t *testing.T) {
	t.Run("should get the ID of a cell", func(t *testing.T) {
		mockDb, mock, _ := sqlmock.New()
		dialector := postgres.New(postgres.Config{
			Conn:       mockDb,
			DriverName: "postgres",
		})

		mock.ExpectQuery(`SELECT \* FROM .+ WHERE id = \$1`).WithArgs("1").WillReturnRows(sqlmock.NewRows([]string{"id", "raw_value", "row_index", "column_index", "spreadsheet_id"}).AddRow(1, "Test Cell", 1, 1, 1))
		mock.ExpectQuery(`SELECT \* FROM .+ WHERE id = \$1`).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "raw_value", "row_index", "column_index", "spreadsheet_id"}).AddRow(1, "Test Cell", 1, 1, 1))

		db, _ := gorm.Open(dialector, &gorm.Config{})
		customCtx := &common.CustomContext{
			Database: db,
		}
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))
		ctx := common.CreateContext(customCtx, srv)

		gql := client.New(ctx)
		resp := struct {
			GetCell struct {
				ID string
			}
		}{}

		q := `query getCellID {
			getCell(id: "1") {
				id
			}
		}`

		gql.MustPost(q, &resp)

		require.NotNil(t, resp.GetCell)
		require.Equal(t, "1", resp.GetCell.ID)
	})

}

func TestCellResolver_RawValue(t *testing.T) {
	t.Run("should get the raw value of a cell", func(t *testing.T) {
		mockDb, mock, _ := sqlmock.New()
		dialector := postgres.New(postgres.Config{
			Conn:       mockDb,
			DriverName: "postgres",
		})

		mock.ExpectQuery(`SELECT \* FROM .+ WHERE id = \$1`).WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "raw_value", "row_index", "column_index", "spreadsheet_id"}).AddRow(1, "Test Cell", 1, 1, 1))

		db, _ := gorm.Open(dialector, &gorm.Config{})
		customCtx := &common.CustomContext{
			Database: db,
		}
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))
		ctx := common.CreateContext(customCtx, srv)

		gql := client.New(ctx)
		resp := struct {
			GetCell struct {
				RawValue string
			}
		}{}

		q := `query getCellRawValue {
			getCell(id: "1") {
				rawValue
			}
		}`

		gql.MustPost(q, &resp)

		require.NotNil(t, resp.GetCell)
		require.Equal(t, "Test Cell", resp.GetCell.RawValue)
	})

}
func TestQueryResolver_GetCellsBySpreadsheetID(t *testing.T) {
	t.Run("should return cells for a given spreadsheet ID", func(t *testing.T) {
		// Mock the database and prepare expectations
		mockDB, mock, _ := sqlmock.New()
		dialector := postgres.New(postgres.Config{
			Conn:       mockDB,
			DriverName: "postgres",
		})
		mock.ExpectQuery(`SELECT \* FROM .+ WHERE \(spreadsheet_id = \$1`).WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "raw_value", "row_index", "column_index", "spreadsheet_id"}).
				AddRow(1, "Test Cell 1", 0, 0, 1).
				AddRow(2, "Test Cell 2", 0, 1, 1).
				AddRow(3, "Test Cell 3", 1, 0, 1))

		// Create a test context with the mocked database
		db, _ := gorm.Open(dialector, &gorm.Config{})
		customCtx := &common.CustomContext{
			Database: db,
		}
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))
		ctx := common.CreateContext(customCtx, srv)

		gql := client.New(ctx)
		resp := struct {
			GetCellsBySpreadsheetId []struct {
				RawValue string
			}
		}{}

		q := `query getCellsBySpreadsheetId {
			getCellsBySpreadsheetId(spreadsheetId: "1") {
				rawValue
			}
		}`
		gql.MustPost(q, &resp)

		assert.NotNil(t, resp.GetCellsBySpreadsheetId)
		assert.Equal(t, 3, len(resp.GetCellsBySpreadsheetId))
		assert.Equal(t, "Test Cell 1", resp.GetCellsBySpreadsheetId[0].RawValue)
		assert.Equal(t, "Test Cell 2", resp.GetCellsBySpreadsheetId[1].RawValue)
		assert.Equal(t, "Test Cell 3", resp.GetCellsBySpreadsheetId[2].RawValue)

	})
}
func TestQueryResolver_GetCell(t *testing.T) {
	t.Run("should return a cell for a given ID", func(t *testing.T) {
		// Mock the database and prepare expectations
		mockDB, mock, _ := sqlmock.New()
		dialector := postgres.New(postgres.Config{
			Conn:       mockDB,
			DriverName: "postgres",
		})
		mock.ExpectQuery(`SELECT \* FROM .+ WHERE id = \$1`).WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "raw_value", "row_index", "column_index", "spreadsheet_id"}).
				AddRow(1, "Test Cell", 0, 0, 1))

		// Create a test context with the mocked database
		db, _ := gorm.Open(dialector, &gorm.Config{})
		customCtx := &common.CustomContext{
			Database: db,
		}
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))
		ctx := common.CreateContext(customCtx, srv)

		gql := client.New(ctx)
		resp := struct {
			GetCell struct {
				ID       string
				RawValue string
			}
		}{}

		q := `query getCell {
			getCell(id: "1") {
				id
				rawValue
			}
		}`
		gql.MustPost(q, &resp)

		assert.NotNil(t, resp.GetCell)
		assert.Equal(t, "1", resp.GetCell.ID)
		assert.Equal(t, "Test Cell", resp.GetCell.RawValue)
	})
}

func TestMutationResolver_CreateCell(t *testing.T) {
	t.Run("should create a new cell", func(t *testing.T) {
		// Mock the database and prepare expectations
		mockDB, mock, _ := sqlmock.New()
		dialector := postgres.New(postgres.Config{
			Conn:       mockDB,
			DriverName: "postgres",
		})
		mock.ExpectQuery(`SELECT \* FROM .+ WHERE id = \$1`).WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "row_count", "column_count"}).AddRow(1, 1, 1))

		// Create a test context with the mocked database
		db, _ := gorm.Open(dialector, &gorm.Config{})
		customCtx := &common.CustomContext{
			Database: db,
		}
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))
		ctx := common.CreateContext(customCtx, srv)

		gql := client.New(ctx)
		resp := struct {
			CreateCell struct {
				ID string
			}
		}{}

		q := `mutation createCell {
			createCell(input: {
				rawValue: "Test Cell",
				rowIndex: 0,
				columnIndex: 0,
				spreadsheetId: "1"
			}) {
				id
			}
		}`

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "cells"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "Test Cell", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()
		// expect panic here

		gql.MustPost(q, &resp)

		assert.NotNil(t, resp.CreateCell)
		assert.Equal(t, "1", resp.CreateCell.ID)
	})

	t.Run("should fail to create a cell outside column and row counts", func(t *testing.T) {
		// Mock the database and prepare expectations
		mockDB, mock, _ := sqlmock.New()
		dialector := postgres.New(postgres.Config{
			Conn:       mockDB,
			DriverName: "postgres",
		})
		mock.ExpectQuery(`SELECT \* FROM .+ WHERE id = \$1`).WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "row_count", "column_count"}).AddRow(1, 1, 1))

		// Create a test context with the mocked database
		db, _ := gorm.Open(dialector, &gorm.Config{})
		customCtx := &common.CustomContext{
			Database: db,
		}
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))
		ctx := common.CreateContext(customCtx, srv)

		gql := client.New(ctx)
		resp := struct {
			CreateCell *model.Cell
		}{}

		q := `mutation createCell {
			createCell(input: {
				rawValue: "Test Cell",
				rowIndex: 2,
				columnIndex: 2,
				spreadsheetId: "1"
			}) {
				id
			}
		}`
		defer func() {
			r := recover()
			assert.NotNil(t, r, "panic should have occurred")
		}()

		gql.MustPost(q, &resp)
		assert.Nil(t, resp.CreateCell)
	})
}
