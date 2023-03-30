resource "materialize_materialized_view" "simple_materialized_view" {
  name          = "simple_materialized_view"
  schema_name   = materialize_schema.schema.name
  database_name = materialize_database.database.name

  statement = <<SQL
SELECT
    *
FROM
    ${materialize_table.simple_table.qualified_name}
SQL

  depends_on = [materialize_table.simple_table]
}

resource "materialize_materialized_view" "simple_materialized_view" {
  name          = "simple_materialized_view"
  schema_name   = materialize_schema.schema.name
  database_name = materialize_database.database.name

  statement = "SELECT * FROM materialize.public.simple_table"
}
