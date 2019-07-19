package main

import (
	"fmt"
	"github.com/steinfletcher/t8/config"
	"github.com/steinfletcher/t8/git"
	"github.com/steinfletcher/t8/io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	data, err := ioutil.ReadFile("./config/testdata/t8.hcl")
	if err != nil {
		panic(err)
	}

	conf, err := config.New(data)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	prompt := io.NewPrompt(io.NewReader())

	err = prompt.ReadParameters(&conf)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = git.Clone("https://github.com/steinfletcher/apitest", conf.TargetDir)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	rmCmd := exec.Command("rm", "-rf", fmt.Sprintf("%s/.git", conf.TargetDir))
	if err = rmCmd.Run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = filepath.Walk(conf.TargetDir, func(path string, f os.FileInfo, err error) error {
		fmt.Println(path)
		return nil
	})

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
