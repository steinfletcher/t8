package io

import (
	"github.com/pkg/errors"
	"github.com/steinfletcher/t8"
	"strings"
)

type Prompt interface {
	ReadParameters(config *t8.Config) error
}

type prompt struct {
	reader Reader
}

func NewPrompt(reader Reader) Prompt {
	return prompt{reader: reader}
}

func (p prompt) ReadParameters(config *t8.Config) error {
	for i, parameter := range config.Parameters {
		result, err := p.reader.Read(parameter)
		if err != nil {
			return errors.Wrap(err, "failed to read user input")
		}

		if parameter.Type == t8.List {
			config.Parameters[i].Actual = strings.Split(result, ",")
			continue
		}
		config.Parameters[i].Actual = result
	}
	return nil
}

