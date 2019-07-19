package git_test

import (
	"github.com/steinfletcher/t8/git"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestClone(t *testing.T) {
	cases := map[string]struct {
		url             string
		targetDir       string
		expectedCommand string
	}{
		"clones with URL and targetDir": {
			url:             "https://github.com/org/repo",
			targetDir:       "myRepo",
			expectedCommand: "git clone --quiet --depth 1 https://github.com/org/repo myRepo",
		},
		"clones with repo name targetDir if not provided": {
			url:             "https://github.com/org/repo",
			expectedCommand: "git clone --quiet --depth 1 https://github.com/org/repo ",
		},
		"clones with org/repo": {
			url:             "org/repo",
			expectedCommand: "git clone --quiet --depth 1 git@github.com:org/repo ",
		},
	}
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			mockExecutor := &executorMock{}
			fetch := git.NewFetchConfig(mockExecutor)

			err := fetch(test.url, test.targetDir)

			assert.NoError(t, err)
			assert.Equal(t, test.expectedCommand, mockExecutor.capturedCommand)
		})
	}
}

type executorMock struct {
	capturedCommand string
}

func (e *executorMock) Exec(name string, arg ...string) error {
	args := append([]string{name}, arg...)
	e.capturedCommand = strings.Join(args, " ")
	return nil
}
