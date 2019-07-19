package git

import (
	"fmt"
	"github.com/steinfletcher/t8/domain"
	"github.com/steinfletcher/t8/shell"
	"strings"
)

// NewFetchConfig is an implementation of FetchTemplate that fetches the template using git
func NewFetchConfig(executor shell.Executor) domain.FetchTemplate {
	return func(url, targetDir string) error {
		if isOrgAndRepo(url) {
			parts := strings.Split(url, "/")
			url = fmt.Sprintf("git@github.com:%s/%s", parts[0], parts[1])
		}
		return executor.Exec("git", "clone", "--quiet", "--depth", "1", url, targetDir)
	}
}

func isOrgAndRepo(url string) bool {
	return !strings.HasPrefix(url, "git@") &&
		!strings.HasPrefix(url, "http") &&
		len(strings.Split(url, "/")) == 2
}
