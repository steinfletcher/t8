name = "My amazing app"
description = "Template to create an amazing application"

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
