package main

import (
	"github.com/steinfletcher/t8/cmd"
	"github.com/steinfletcher/t8/git"
	"github.com/steinfletcher/t8/shell"
	"github.com/steinfletcher/t8/template"
	"log"
	"os"
)

var version string

func main() {
	exec := shell.NewExecutor()
	fetchConfig := git.NewFetchConfig(exec)
	prompt := shell.NewStdInReader()
	templateRenderer := template.NewTemplateRenderer()

	if err := cmd.Run(fetchConfig, prompt, templateRenderer, os.Args, version); err != nil {
		log.Fatal(err)
	}
}
