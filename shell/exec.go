package shell

import "os/exec"

type Executor interface {
	Exec(name string, arg ...string) error
}

type shellExecutor struct{}

func NewExecutor() Executor {
	return &shellExecutor{}
}

func (r *shellExecutor) Exec(name string, args ...string) error {
	return exec.Command(name, args...).Run()
}
