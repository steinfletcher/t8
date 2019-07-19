package io_test

import (
	"fmt"
	"github.com/steinfletcher/t8"
	"github.com/steinfletcher/t8/io"
	"github.com/stretchr/testify/assert"
	"testing"
)

var conf = t8.Config{
	Name: "my project",
	Parameters: []t8.Parameter{
		{
			Name: "go_version",
			Type: t8.String,
		},
		{
			Name: "authors",
			Type: t8.List,
		},
		{
			Name: "sql_dialect",
			Type: t8.Option,
		},
	},
}

func TestPrompt_ReadParameters(t *testing.T) {
	prompt := io.NewPrompt(mockReader{})

	err := prompt.ReadParameters(&conf)

	assert.NoError(t, err)
	assert.Equal(t, "1.12", conf.GetParameter("go_version").Actual)
	assert.Equal(t, []string{"yuki", "mei"}, conf.GetParameter("authors").Actual)
	assert.Equal(t, "postgresql", conf.GetParameter("sql_dialect").Actual)
}

type mockReader struct {
	t *testing.T
}

func (m mockReader) Read(parameter t8.Parameter) (string, error) {
	if parameter.Name == "go_version" {
		return "1.12", nil
	}

	if parameter.Name == "authors" {
		return "yuki,mei", nil
	}

	if parameter.Name == "sql_dialect" {
		return "postgresql", nil
	}

	return "", fmt.Errorf("parameter name not expected: '%v'", parameter.Name)
}
