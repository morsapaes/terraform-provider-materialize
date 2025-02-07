package datasources

import (
	"context"
	"testing"

	"github.com/MaterializeInc/terraform-provider-materialize/pkg/testhelpers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestIndexDatasource(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"schema_name":   "schema",
		"database_name": "database",
	}
	d := schema.TestResourceDataRaw(t, Index().Schema, in)
	r.NotNil(d)

	testhelpers.WithMockDb(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		ir := mock.NewRows([]string{"id", "index_name", "obj_name", "obj_schema_name", "obj_database_name"}).
			AddRow("u1", "index", "obj", "schema", "database")
		mock.ExpectQuery(`
			SELECT
				mz_indexes.id,
				mz_indexes.name AS index_name,
				mz_objects.name AS obj_name,
				mz_schemas.name AS obj_schema_name,
				mz_databases.name AS obj_database_name
			FROM mz_indexes
			JOIN mz_objects
				ON mz_indexes.on_id = mz_objects.id
			LEFT JOIN mz_schemas
				ON mz_objects.schema_id = mz_schemas.id
			LEFT JOIN mz_databases
				ON mz_schemas.database_id = mz_databases.id
			WHERE mz_databases.name = 'database'
			AND mz_objects.type IN \('source', 'view', 'materialized-view'\)
			AND mz_schemas.name = 'schema';
		`).WillReturnRows(ir)

		if err := indexRead(context.TODO(), d, db); err != nil {
			t.Fatal(err)
		}
	})

}
