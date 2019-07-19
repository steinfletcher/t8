package git

import (
	"os"
	"os/exec"
)

func Clone(url, target string) error {
	var err error
	cmd := exec.Command("git", "clone", "--depth", "1", url, target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return err
	}
	return nil
}
