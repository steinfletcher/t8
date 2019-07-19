package template

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func Render(repoPath string) error {
	return filepath.Walk(repoPath, func(path string, f os.FileInfo, err error) error {
		fmt.Println(path)
		return errors.New("not yet implemented")
	})
}
