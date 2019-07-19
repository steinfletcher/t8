*This is experimental and doesn't work yet - do not use*

# t8

`t8` (pronounced template) is a simple CLI application that renders templates defined on Github (and other locations).

Inspired by [giter8](http://www.foundweekends.org/giter8/).

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
$ t8 https://github.com/org/spring-boot-template.t8 my-amazing-app

Enter the required parameters to generate "Spring boot app".

project name[acme]: My amazing app
go version[1.12]: 

Template created: /home/stein/code/my-amazing-app
```

### Pass input parameters

```bash
$ t8 -project_name=acme -go_version=1.12 https://github.com/org/go-echo-template.t8 my-amazing-app

Template created: /home/stein/code/my-amazing-app
```

## Roadmap

* Render templates defined on local filesystem 
* Render templates defined at a remote http location, e.g. github
* Support more complex parameters, e.g. maps and objects
* Support var file to initialise template from ala terraform (is there a use-case for this?)
* Search templates. Should be able to build an index of templates from github - find repos with t8.hcl/yml in root of repo, or repos ending in `.t8`
