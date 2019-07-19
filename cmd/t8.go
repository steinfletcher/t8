package cmd

import (
	"fmt"
	"github.com/steinfletcher/t8/domain"
	"github.com/steinfletcher/t8/hcl"
	"github.com/steinfletcher/t8/shell"
	"io/ioutil"
	"log"
	"os"
)

func Run(
	fetchConfig domain.FetchTemplate,
	prompt domain.PromptReader,
	template domain.TemplateRenderer,
	args []string,
) {
	command, err := shell.ParseArgs(args)
	exitOnErr(err)

	if f, ok := command.GetFlag("Debug"); ok && f.Value == "true" {
		domain.IsDebug = true
	}

	if command.Name == "help" {
		printHelp()
		os.Exit(0)
	}

	if command.Name == "version" {
		fmt.Println("show version...")
		os.Exit(0)
	}

	if command.Name == "new" {
		runNewCommand(command, prompt, fetchConfig, template)
		os.Exit(0)
	}

	printHelp()
	os.Exit(1)
}

func runNewCommand(command domain.Command, prompt domain.PromptReader, fetchConfig domain.FetchTemplate, template domain.TemplateRenderer) {
	if len(command.Args) == 0 {
		fmt.Println("Please specify a template location")
		printHelp()
		os.Exit(1)
	}

	repo := command.Args[0]
	targetDir := "repo"
	if len(command.Args) > 1 {
		targetDir = command.Args[1]
	}

	domain.PrepareTargetDir(prompt, targetDir)
	err := domain.PrepareConfig(fetchConfig, repo, targetDir)
	exitOnErr(err)

	printFileToConsole("before.t8", targetDir)

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/t8.yml", targetDir))
	if err != nil {
		data, err = ioutil.ReadFile(fmt.Sprintf("%s/t8.hcl", targetDir))
		exitOnErr(err)
	}

	config, err := hcl.ParseConfig(data)
	exitOnErr(err)

	validConfig, err := domain.PromptForRequiredParameters(prompt, config, command)
	exitOnErr(err)

	if domain.IsDebug {
		log.Printf("Final config: %+v\n", validConfig)
	}

	template.Render(targetDir, validConfig)

	printFileToConsole("after.t8", targetDir)

	fmt.Printf("\nTemplate created: '%v'\n\n", targetDir)
}

func printFileToConsole(f, targetDir string) {
	filePath := fmt.Sprintf("%s/%s", targetDir, f)
	if _, err := os.Stat(filePath); err == nil {
		data, err := ioutil.ReadFile(filePath)
		exitOnErr(err)
		fmt.Println(string(data))
	}
}

func exitOnErr(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`NAME:
   t8 - a customizable template generator

USAGE:
   t8 new myorg/myrepo.t8 localDir -ProjectName=my-project -GoVersion=1.12

VERSION:
   0.1.1

COMMANDS:
     new         Render a new template
     help,    -h Shows help
     version, -v Shows version`)
}
