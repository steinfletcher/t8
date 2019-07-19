name = "My Amazing App"

parameter "ProjectName" {
  type = "string"
  description = "the project name"
  default = "acme"
}

parameter "Authors" {
  type = "option"
  description = "the project authors"
  default = [
    "yuki@foo.com",
    "mei@bar.com",
  ]
}

parameter "SqlDialect" {
  type = "option"
  description = "the SQL dialect"
  default = [
    "postgresql",
    "mysql",
  ]
}

excludePath "Postgres" {
  paths = [
    "^/postgres/.*$"
  ]
  parameterName = "SqlDialect"
  operator = "notEqual"
  parameterValue = "postgresql"
}

excludePath "Static" {
  paths = [
    "test.sh"
  ]
}
