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

var inConfluentSchemaRegistry = map[string]interface{}{
	"name":                      "conn",
	"schema_name":               "schema",
	"database_name":             "database",
	"service_name":              "service",
	"url":                       "http://localhost:8081",
	"ssl_certificate_authority": []interface{}{map[string]interface{}{"secret": []interface{}{map[string]interface{}{"name": "ssl"}}}},
	"ssl_certificate":           []interface{}{map[string]interface{}{"secret": []interface{}{map[string]interface{}{"name": "ssl"}}}},
	"ssl_key":                   []interface{}{map[string]interface{}{"name": "ssl"}},
	"password":                  []interface{}{map[string]interface{}{"name": "password"}},
	"username":                  []interface{}{map[string]interface{}{"text": "user"}},
	"ssh_tunnel":                []interface{}{map[string]interface{}{"name": "tunnel"}},
	"aws_privatelink":           []interface{}{map[string]interface{}{"name": "privatelink"}},
}

func TestResourceConnectionConfluentSchemaRegistryCreate(t *testing.T) {
	r := require.New(t)
	d := schema.TestResourceDataRaw(t, ConnectionConfluentSchemaRegistry().Schema, inConfluentSchemaRegistry)
	r.NotNil(d)

	testhelpers.WithMockDb(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		// Create
		mock.ExpectExec(
			`CREATE CONNECTION "database"."schema"."conn" TO CONFLUENT SCHEMA REGISTRY \(URL 'http://localhost:8081', USERNAME = 'user', PASSWORD = SECRET "database"."schema"."password", SSL CERTIFICATE AUTHORITY = SECRET "database"."schema"."ssl", SSL CERTIFICATE = SECRET "database"."schema"."ssl", SSL KEY = SECRET "database"."schema"."ssl", AWS PRIVATELINK "database"."schema"."privatelink", SSH TUNNEL "database"."schema"."tunnel"\)`,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		// Query Id
		ir := mock.NewRows([]string{"id"}).AddRow("u1")
		mock.ExpectQuery(`
			SELECT
				mz_connections.id,
				mz_connections.name AS connection_name,
				mz_schemas.name AS schema_name,
				mz_databases.name AS database_name,
				mz_connections.type AS connection_type
			FROM mz_connections
			JOIN mz_schemas
				ON mz_connections.schema_id = mz_schemas.id
			JOIN mz_databases
				ON mz_schemas.database_id = mz_databases.id
			WHERE mz_connections.name = 'conn'
			AND mz_databases.name = 'database'
			AND mz_schemas.name = 'schema';`).WillReturnRows(ir)

		// Query Params
		ip := sqlmock.NewRows([]string{"connection_name", "schema_name", "database_name"}).
			AddRow("conn", "schema", "database")
		mock.ExpectQuery(readConnection).WillReturnRows(ip)

		if err := connectionConfluentSchemaRegistryCreate(context.TODO(), d, db); err != nil {
			t.Fatal(err)
		}
	})

}
