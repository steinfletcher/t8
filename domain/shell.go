package domain

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
)

//go:generate mockgen -source=shell.go -destination=../mocks/shell.go -package=mocks

// PromptReader is an interface for reading user input
type PromptReader interface {
	String(parameter Parameter) (string, error)
	Options(parameter Parameter) (string, error)
	Confirm(question string) (bool, error)
}

// Command represents the result of a shell command,
// e.g. "t8 <Name> <Flags> <Args>" might be "t8 new -Key=Val https://github.com/org/repo dir"
type Command struct {
	Name  string
	Flags []Flag
	Args  []string
}

// Flag models a command line flag, e.g. "-Key=ParameterValue"
type Flag struct {
	Key   string
	Value string
}

// GetFlag gets the command line flag by name
func (c Command) GetFlag(key string) (Flag, bool) {
	for _, flag := range c.Flags {
		if flag.Key == key {
			return flag, true
		}
	}
	return Flag{}, false
}

// PromptForRequiredParameters reads the config and prompts the user for any required parameters that do not have a
// value set. Values can be set already by passing flags via the command line
func PromptForRequiredParameters(prompt PromptReader, config Config, command Command) (Config, error) {
	for i, parameter := range config.Parameters {
		if flag, ok := command.GetFlag(parameter.Name); ok {
			config.Parameters[i].Actual = flag.Value
			continue
		}

		if parameter.Type == Option {
			result, err := prompt.Options(parameter)
			if err != nil {
				return Config{}, errors.Wrap(err, "failed to read user input")
			}
			config.Parameters[i].Actual = result
			continue
		}

		result, err := prompt.String(parameter)
		if err != nil {
			return Config{}, errors.Wrap(err, "failed to read user input")
		}
		config.Parameters[i].Actual = result
	}
	return config, nil
}

// PrepareTargetDir prompts user if the targetDir is not empty and optionally deletes the targetDir it if the user
// consents
func PrepareTargetDir(prompt PromptReader, targetDir string) {
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		shouldDelete, err := prompt.Confirm(fmt.Sprintf("Directory '%s' is not empty, overwrite it?", targetDir))
		exitOnErr(err)

		if !shouldDelete {
			log.Println("Exiting.")
			os.Exit(0)
		}

		err = os.RemoveAll(targetDir)
		log.Println(fmt.Sprintf("Deleted directory '%s'.", targetDir))
		exitOnErr(err)
	}
}

func exitOnErr(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
