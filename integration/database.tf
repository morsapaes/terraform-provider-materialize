resource "materialize_database" "database" {
  name = "example"
}

data "materialize_database" "all" {}

data "materialize_current_database" "default" {}
