package shell

import (
	"github.com/manifoldco/promptui"
	"github.com/steinfletcher/t8/domain"
	"strings"
)

// NewStdInReader creates a new input PromptReader to read from standard input
func NewStdInReader() domain.PromptReader {
	return stdInReader{}
}

// stdInReader wraps promptui
type stdInReader struct{}

func (stdInReader) String(parameter domain.Parameter) (string, error) {
	prompt := promptui.Prompt{
		Validate: func(s string) error {
			for _, validate := range parameter.Validators {
				if err := validate(s); err != nil {
					return err
				}
			}
			return nil
		},
		Label: parameter.Name + " ",
	}

	if parameter.Default != nil {
		prompt.Default = parameter.DefaultString()
	}

	return prompt.Run()
}

func (stdInReader) Options(parameter domain.Parameter) (string, error) {
	prompt := promptui.Select{
		Label: parameter.Name,
		Items: parameter.Default.([]string),
	}
	_, result, err := prompt.Run()
	return result, err
}

func (stdInReader) Confirm(question string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     question,
		IsConfirm: true,
		Default:   "N",
	}

	result, err := prompt.Run()
	if err == promptui.ErrAbort {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if result == "y" {
		return true, nil
	}

	return false, nil
}

func ParseArgs(input []string) (domain.Command, error) {
	cmd := domain.Command{}
	var args []string
	for _, field := range input {
		// parse flags
		if strings.HasPrefix(field, "-") {
			token := strings.TrimLeft(field, "-")
			parts := strings.Split(token, "=")

			// treat -h as the "help" command
			if len(parts) == 1 && parts[0] == "h" {
				cmd.Name = "help"
				return cmd, nil
			}

			// treat -v as the "version" command
			if len(parts) == 1 && parts[0] == "v" {
				cmd.Name = "version"
				return cmd, nil
			}

			if len(parts) != 2 {
				return domain.Command{}, domain.ErrInvalidArgs("expected flag to take the form 'key=value'")
			}
			cmd.Flags = append(cmd.Flags, domain.Flag{Key: parts[0], Value: parts[1]})
			continue
		}

		args = append(args, field)
	}

	if len(args) < 2 {
		return domain.Command{}, domain.ErrInvalidArgs("expected a command")
	}

	// add command name
	cmd.Name = args[1]

	// rest of the args are considered as trailing args
	if len(args) > 2 {
		cmd.Args = args[2:]
	}

	return cmd, nil
}
