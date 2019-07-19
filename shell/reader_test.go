package shell_test

import (
	"github.com/steinfletcher/t8/domain"
	"github.com/steinfletcher/t8/shell"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	cases := map[string]struct {
		input   string
		command domain.Command
	}{
		"command only": {
			input:   "t8 help",
			command: domain.Command{Name: "help"},
		},
		"command with mixed flags and args": {
			input: "t8 new org/repo myDir --ProjectName=wowzers -GoVersion=1.12",
			command: domain.Command{
				Name: "new",
				Flags: []domain.Flag{
					{
						Key: "ProjectName", Value: "wowzers",
					},
					{
						Key: "GoVersion", Value: "1.12",
					},
				},
				Args: []string{"org/repo", "myDir"},
			},
		},
		"command with flags only": {
			input: "t8 new --ProjectName=wowzers -GoVersion=1.12",
			command: domain.Command{
				Name: "new",
				Flags: []domain.Flag{
					{
						Key: "ProjectName", Value: "wowzers",
					},
					{
						Key: "GoVersion", Value: "1.12",
					},
				},
			},
		},
		"command with args only": {
			input: "t8 new org/repo",
			command: domain.Command{
				Name: "new",
				Args: []string{"org/repo"},
			},
		},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			command, err := shell.ParseArgs(strings.Fields(testCase.input))

			assert.NoError(t, err)
			assert.Equal(t, testCase.command, command)
		})
	}
}
