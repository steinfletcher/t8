package t8_test

import (
	"github.com/pkg/errors"
	"github.com/steinfletcher/t8"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var config = t8.Config{
	Name: "my project",
	Parameters: []t8.Parameter{
		{
			Name:    "myName",
			Default: "myProject",
		},
	},
}

func TestConfig_GetParameter(t *testing.T) {
	cases := map[string]struct {
		parameterName     string
		expectedParameter *t8.Parameter
	}{
		"success": {
			parameterName:     "myName",
			expectedParameter: &t8.Parameter{Default: "myProject", Name: "myName"},
		},
		"not found": {
			parameterName:     "doesn't Exist",
			expectedParameter: nil,
		},
	}
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			parameter := config.GetParameter(test.parameterName)

			assert.Equal(t, test.expectedParameter, parameter)
		})
	}
}

func TestComposeValidators(t *testing.T) {
	helloValidator := func(in string) error {
		if strings.Contains(in, "hello") {
			return nil
		}
		return errors.New("should contain hello")
	}

	worldValidator := func(in string) error {
		if strings.Contains(in, "world") {
			return nil
		}
		return errors.New("should contain world")
	}

	t8.ComposeValidators(helloValidator, worldValidator)
}

func TestComposeValidators2(t *testing.T) {
	helloValidator := func(in string) error {
		if strings.Contains(in, "hello") {
			return nil
		}
		return errors.New("should contain hello")
	}

	worldValidator := func(in string) error {
		if strings.Contains(in, "world") {
			return nil
		}
		return errors.New("should contain world")
	}

	cases := map[string]struct {
		input       string
		expectedErr error
	}{
		"passes both validators": {
			input:       "hello world",
			expectedErr: nil,
		},
		"validator1 fail": {
			input:       "hello",
			expectedErr: errors.New("should contain world"),
		},
		"validator2 fail": {
			input:       "world",
			expectedErr: errors.New("should contain hello"),
		},
	}
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			validate := t8.ComposeValidators(helloValidator, worldValidator)

			err := validate(test.input)

			assert.Equal(t, test.expectedErr, err)
		})
	}
}
