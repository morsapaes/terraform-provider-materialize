package resources

import (
	"context"
	"testing"

	"github.com/MaterializeInc/terraform-provider-materialize/pkg/testhelpers"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestResourceDatabaseCreate(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name": "database",
	}
	d := schema.TestResourceDataRaw(t, Database().Schema, in)
	r.NotNil(d)

	testhelpers.WithMockDb(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		// Create
		mock.ExpectExec(
			`CREATE DATABASE "database";`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		// Query Id
		ir := mock.NewRows([]string{"id", "database_name"}).AddRow("u1", "database")
		mock.ExpectQuery(`SELECT id, name AS database_name FROM mz_databases WHERE name = 'database'`).WillReturnRows(ir)

		// Query Params
		ip := mock.NewRows([]string{"id", "database_name"}).AddRow("u1", "database")
		mock.ExpectQuery(`SELECT id, name AS database_name FROM mz_databases WHERE id = 'u1';`).WillReturnRows(ip)

		if err := databaseCreate(context.TODO(), d, db); err != nil {
			t.Fatal(err)
		}
	})

}

func TestResourceDatabaseDelete(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"name": "database",
	}
	d := schema.TestResourceDataRaw(t, Database().Schema, in)
	r.NotNil(d)

	testhelpers.WithMockDb(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(`DROP DATABASE "database";`).WillReturnResult(sqlmock.NewResult(1, 1))

		if err := databaseDelete(context.TODO(), d, db); err != nil {
			t.Fatal(err)
		}
	})

}
