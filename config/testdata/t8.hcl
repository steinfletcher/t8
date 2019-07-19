name = "My Amazing App"

parameter "project_name" {
  type = "string"
  description = "the project name"
  default = "acme"
}

parameter "authors" {
  type = "list"
  description = "the project authors"
  default = [
    "yuki@foo.com",
    "mei@bar.com",
  ]
}

parameter "sql_dialect" {
  type = "option"
  description = "the SQL dialect"
  default = [
    "postgresql",
    "mysql",
  ]
}
