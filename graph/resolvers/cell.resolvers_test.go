package resolvers

import (
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
	"github.com/vijaykramesh/gql-sheets/graph/common"
	"github.com/vijaykramesh/gql-sheets/graph/generated"
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
