package io

import (
	"github.com/manifoldco/promptui"
	"github.com/steinfletcher/t8"
)

// NewReader creates a new input Reader
func NewReader() Reader {
	return stdInReader{}
}

// stdInReader wraps promptui so it can be mocked in tests
type stdInReader struct{}

// Reader is used to read user input for the given parameter
type Reader interface {
	Read(parameter t8.Parameter) (string, error)
}

// Read reads the parameter spec from user input using promptui
func (stdInReader) Read(parameter t8.Parameter) (string, error) {
	if parameter.Type == t8.Option {
		prompt := promptui.Select{
			Label: parameter.Name,
			Items: parameter.Default.([]interface{}),
		}
		_, result, err := prompt.Run()
		return result, err
	}

	prompt := promptui.Prompt{
		Validate: func(s string) error {
			for _, validate := range parameter.Validators {
				if err := validate(s); err != nil {
					return err
				}
			}
			return nil
		},
		Label:     parameter.Name + " ",
	}
	return prompt.Run()
}
