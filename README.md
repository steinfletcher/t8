[![Build Status](https://travis-ci.com/steinfletcher/t8.svg?token=iwqcySR3NwvVKvZeyk52&branch=master)](https://travis-ci.com/steinfletcher/t8)

*This is experimental and the API is likely to change*

# t8

`t8` (pronounced template) is a simple CLI application that renders templates defined on Github (and other locations).

Inspired by [giter8](http://www.foundweekends.org/giter8/).

## Install

Install using Go

```bash
go get github.com/steinfletcher/t8/cmd/t8
```

Or download a prebuilt binary from github releases.

## Use cases

* A scaffolding tool for generating boilerplate applications (like Yeoman or sbt minus the build step)
* Automate generation of config files

## Features

* Uses Go templating
* Define config as HCL or YAML
* Interactive CLI to prompt for parameters
* Pass input parameters as arguments (useful in CI)

## Usage

## Interactive CLI

```bash
$ t8 new https://github.com/myOrg/myTemplate my-amazing-app

Enter the required parameters to generate "My Amazing App".

project name[acme]: My amazing app
go version[1.12]: 

Template created: /home/stein/code/my-amazing-app
```

### Pass input parameters

```bash
$ t8 -ProjectName=acme -GoVersion=1.12 https://github.com/org/go-echo-template.t8 my-amazing-app

Template created: /home/stein/code/my-amazing-app
```

## Configuration

Create a Go template and host it on github. Create a file called `t8.hcl` or `t8.yml` in the root of the project. This is the configuration file where you can configure your generator.

### Parameters

A `parameter` is a variable set at runtime and is made accessible to your Go template. You can also define defaults

```hcl
parameter "ProjectName" {
  type = "string"
  description = "the project name"
  default = "acme"
}
```

In this example the user will be prompted to enter the project name or provide it as a command line flag. If the user does not define this variable the default value is used.

### Exclude Paths

You can exclude the generation of files and directories using the `excludePath` variable. 

```hcl
excludePath "Scripts" {
  paths = [
    "test.sh",
  ]
}
```

This configures an unconditional exclusion on a path pattern - the `test.sh` file will be excluded from the final generated content. You can also exclude paths based on the value of a parameter.

```hcl
excludePath "Postgres" {
  paths = [
    "^/postgres/.*$"
  ]
  parameterName = "SqlDialect"
  operator = "notEqual"
  parameterValue = "postgresql"
}
```

In this example, the `/postgres` directory will not be generated if the `SqlDialect` parameter is not equal to `postgresql`. The available operators are `equal` and `notEqual`
