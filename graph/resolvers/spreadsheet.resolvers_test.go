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

var (
	loginname = "mrdulin"
	avatarURL = "avatar.jpg"
	score     = 50
	createAt  = "1900-01-01"
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
