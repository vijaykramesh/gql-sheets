package resolvers

import (
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
	"github.com/vijaykramesh/gql-sheets/graph/common"
	"github.com/vijaykramesh/gql-sheets/graph/generated"
	"github.com/vijaykramesh/gql-sheets/graph/model"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestMutationResolver_CreateSpreadsheet(t *testing.T) {
	t.Run("should create a spreadsheet", func(t *testing.T) {
		mockDb, mock, _ := sqlmock.New()
		dialector := postgres.New(postgres.Config{
			Conn:       mockDb,
			DriverName: "postgres",
		})
		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO .+`).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		db, _ := gorm.Open(dialector, &gorm.Config{})
		customCtx := &common.CustomContext{
			Database: db,
		}
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))
		ctx := common.CreateContext(customCtx, srv)

		gql := client.New(ctx)
		resp := struct {
			CreateSpreadsheet *model.Spreadsheet
		}{}

		q := `mutation createSpreadsheet {
			createSpreadsheet(input: {
				name: "Test Spreadsheet"
				rowCount: 10
				columnCount: 5
			}) {
				name
				rowCount
				columnCount
			}
		}`

		gql.MustPost(q, &resp)

		require.NotNil(t, resp.CreateSpreadsheet)
		require.Equal(t, "Test Spreadsheet", resp.CreateSpreadsheet.Name)
		require.Equal(t, 10, resp.CreateSpreadsheet.RowCount)
		require.Equal(t, 5, resp.CreateSpreadsheet.ColumnCount)
	})
}
func TestMutationResolver_UpdateSpreadsheet(t *testing.T) {
	t.Run("should update a spreadsheet", func(t *testing.T) {
		mockDb, mock, _ := sqlmock.New()
		dialector := postgres.New(postgres.Config{
			Conn:       mockDb,
			DriverName: "postgres",
		})

		mock.ExpectQuery(`SELECT \* FROM .+ WHERE id = \$1`).WithArgs("1").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "rowCount", "columnCount"}).AddRow(1, "Test Spreadsheet", 10, 5))
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE .+ SET .+ WHERE .+ "id" = \$\d+`).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		db, _ := gorm.Open(dialector, &gorm.Config{})
		customCtx := &common.CustomContext{
			Database: db,
		}
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))
		ctx := common.CreateContext(customCtx, srv)

		gql := client.New(ctx)
		resp := struct {
			UpdateSpreadsheet *model.Spreadsheet
		}{}

		q := `mutation updateSpreadsheet {
			updateSpreadsheet(id: "1", input: {
				name: "Updated Spreadsheet"
				rowCount: 20
				columnCount: 8
			}) {
				name
				rowCount
				columnCount
			}
		}`

		gql.MustPost(q, &resp)

		require.NotNil(t, resp.UpdateSpreadsheet)
		require.Equal(t, "Updated Spreadsheet", resp.UpdateSpreadsheet.Name)
		require.Equal(t, 20, resp.UpdateSpreadsheet.RowCount)
		require.Equal(t, 8, resp.UpdateSpreadsheet.ColumnCount)
	})
}

func TestQueryResolver_Spreadsheets(t *testing.T) {
	t.Run("should get a list of spreadsheets", func(t *testing.T) {
		mockDb, mock, _ := sqlmock.New()
		dialector := postgres.New(postgres.Config{
			Conn:       mockDb,
			DriverName: "postgres",
		})
		mock.ExpectQuery(`SELECT \* FROM .+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "row_count", "column_count"}).AddRow(1, "Test Spreadsheet 1", 10, 5).AddRow(2, "Test Spreadsheet 2", 20, 8))
		mock.ExpectCommit()

		db, _ := gorm.Open(dialector, &gorm.Config{})
		customCtx := &common.CustomContext{
			Database: db,
		}
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))
		ctx := common.CreateContext(customCtx, srv)

		gql := client.New(ctx)
		resp := struct {
			Spreadsheets []*model.Spreadsheet
		}{}

		q := `query getSpreadsheets {
			spreadsheets {
				name
				rowCount
				columnCount
			}
		}`

		gql.MustPost(q, &resp)

		require.NotNil(t, resp.Spreadsheets)
		require.Equal(t, 2, len(resp.Spreadsheets))

		spreadsheet1 := resp.Spreadsheets[0]
		require.Equal(t, "Test Spreadsheet 1", spreadsheet1.Name)
		require.Equal(t, 10, spreadsheet1.RowCount)
		require.Equal(t, 5, spreadsheet1.ColumnCount)

		spreadsheet2 := resp.Spreadsheets[1]
		require.Equal(t, "Test Spreadsheet 2", spreadsheet2.Name)
		require.Equal(t, 20, spreadsheet2.RowCount)
		require.Equal(t, 8, spreadsheet2.ColumnCount)
	})
}

func TestQueryResolver_GetSpreadsheet(t *testing.T) {
	t.Run("should get a specific spreadsheet", func(t *testing.T) {
		mockDb, mock, _ := sqlmock.New()
		dialector := postgres.New(postgres.Config{
			Conn:       mockDb,
			DriverName: "postgres",
		})

		mock.ExpectQuery(`SELECT \* FROM .+ WHERE id = \$1`).WithArgs("1").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "row_count", "column_count"}).AddRow(1, "Test Spreadsheet", 10, 5))

		db, _ := gorm.Open(dialector, &gorm.Config{})
		customCtx := &common.CustomContext{
			Database: db,
		}
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))
		ctx := common.CreateContext(customCtx, srv)

		gql := client.New(ctx)
		resp := struct {
			GetSpreadsheet *model.Spreadsheet
		}{}

		q := `query getSpreadsheet {
			getSpreadsheet(id: "1") {
				name
				rowCount
				columnCount
			}
		}`

		gql.MustPost(q, &resp)

		require.NotNil(t, resp.GetSpreadsheet)
		require.Equal(t, "Test Spreadsheet", resp.GetSpreadsheet.Name)
		require.Equal(t, 10, resp.GetSpreadsheet.RowCount)
		require.Equal(t, 5, resp.GetSpreadsheet.ColumnCount)
	})
}
func TestSpreadsheetResolver_ID(t *testing.T) {
	t.Run("should get the ID of a spreadsheet", func(t *testing.T) {
		mockDb, mock, _ := sqlmock.New()
		dialector := postgres.New(postgres.Config{
			Conn:       mockDb,
			DriverName: "postgres",
		})

		mock.ExpectQuery(`SELECT \* FROM .+ WHERE id = \$1`).WithArgs("1").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "row_count", "column_count"}).AddRow(1, "Test Spreadsheet", 10, 5))
		// TODO: this seems very weird, why two queries?
		mock.ExpectQuery(`SELECT \* FROM .+ WHERE id = \$1`).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "row_count", "column_count"}).AddRow(1, "Test Spreadsheet", 10, 5))

		db, _ := gorm.Open(dialector, &gorm.Config{})
		customCtx := &common.CustomContext{
			Database: db,
		}
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))
		ctx := common.CreateContext(customCtx, srv)

		gql := client.New(ctx)
		resp := struct {
			GetSpreadsheet struct {
				ID string
			}
		}{}

		q := `query getSpreadsheetID {
			getSpreadsheet(id: "1") {
				id
			}
		}`

		gql.MustPost(q, &resp)

		require.NotNil(t, resp.GetSpreadsheet)
		require.Equal(t, "1", resp.GetSpreadsheet.ID)
	})
}
