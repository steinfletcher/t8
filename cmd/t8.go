package cmd

import (
	"errors"
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
	version string,
) error {
	command, err := shell.ParseArgs(args)
	if err != nil {
		return err
	}

	if f, ok := command.GetFlag("Debug"); ok && f.Value == "true" {
		domain.IsDebug = true
	}

	if command.Name == "help" {
		printHelp(version)
		return nil
	}

	if command.Name == "version" {
		fmt.Println(version)
		return nil
	}

	if command.Name == "new" {
		return runNewCommand(command, prompt, fetchConfig, template)
	}

	return errors.New("expected to receive a command")
}

func runNewCommand(command domain.Command, prompt domain.PromptReader, fetchConfig domain.FetchTemplate, template domain.TemplateRenderer) error {
	if len(command.Args) == 0 {
		return errors.New("expected to receive a template location")
	}

	repo := command.Args[0]
	targetDir := "repo"
	if len(command.Args) > 1 {
		targetDir = command.Args[1]
	}

	domain.PrepareTargetDir(prompt, targetDir)
	err := domain.PrepareConfig(fetchConfig, repo, targetDir)
	if err != nil {
		return err
	}

	err = printFileToConsole("before.t8", targetDir)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/t8.yml", targetDir))
	if err != nil {
		if data, err = ioutil.ReadFile(fmt.Sprintf("%s/t8.hcl", targetDir)); err != nil {
			return err
		}
	}

	config, err := hcl.ParseConfig(data)
	if err != nil {
		return err
	}

	validConfig, err := domain.PromptForRequiredParameters(prompt, config, command)
	if err != nil {
		return err
	}

	if domain.IsDebug {
		log.Printf("Final config: %+v\n", validConfig)
	}

	template.Render(targetDir, validConfig)

	if err := printFileToConsole("after.t8", targetDir); err != nil {
		return err
	}

	fmt.Printf("\nTemplate created: '%v'\n\n", targetDir)
	return nil
}

func printFileToConsole(f, targetDir string) error {
	filePath := fmt.Sprintf("%s/%s", targetDir, f)
	if _, err := os.Stat(filePath); err == nil {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	}
	return nil
}

func printHelp(version string) {
	fmt.Printf(`NAME:
   t8 - a customizable template generator

USAGE:
   t8 new myorg/myrepo.t8 localDir -ProjectName=my-project -GoVersion=1.12

VERSION:
   %s

COMMANDS:
     new         Render a new template
     help,    -h Shows help
     version, -v Shows version\n`, version)
}
