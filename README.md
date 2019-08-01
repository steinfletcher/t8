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

### Interactive CLI

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
