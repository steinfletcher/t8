package domain

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

//go:generate gonum -types=ParameterTypeEnum,RemovePathOperatorEnum

var IsDebug bool

// PrepareConfig fetches the t8 config and removes any git references if they are presents.
func PrepareConfig(fetchConfig FetchTemplate, location, targetDir string) error {
	err := fetchConfig(location, targetDir)
	if err != nil {
		return err
	}

	rmCmd := exec.Command("rm", "-rf", fmt.Sprintf("%s/.git", targetDir))
	return rmCmd.Run()
}

// FetchTemplate fetches the t8 template from the provided source and stores at the given target
type FetchTemplate func(source string, target string) error

// Config is the t8 configuration that is by default defined as hcl or yaml in `t8.hcl` or `t8.yml`
type Config struct {
	Parameters   []Parameter
	ExcludePaths []ExcludePath
	Name         string
}

// ParameterName is a user defined parameter defined in the t8 config. It is passed into the user defined Go template
type Parameter struct {
	Name        string
	Type        ParameterType
	Description string
	Default     interface{}
	Actual      interface{}
	Validators  []Validator
}

// ExcludePath allows the user to configure which files and directories are rendered
type ExcludePath struct {
	ParameterName  string
	ParameterValue string
	Operator       RemovePathOperator
	Paths          []string
}

// ParameterTypeEnum models the different types that may be defined for a parameter
type ParameterTypeEnum struct {
	String string `enum:"string"`
	Option string `enum:"option"`
}

// RemovePathOperatorEnum models the different operations that can be applied when filtering filepaths with parameters
type RemovePathOperatorEnum struct {
	Equal    string `enum:"equal"`
	NotEqual string `enum:"notEqual"`
}

// ErrInvalidArgs is an error that represents a failed t8 user defined command
type ErrInvalidArgs string

func (e ErrInvalidArgs) Error() string {
	return string(e)
}

// GetParameter returns the t8 parameter by name. It returns true if it exists and false if it does not exist
func (c Config) GetParameter(name string) (Parameter, bool) {
	for i, v := range c.Parameters {
		if v.Name == name {
			return c.Parameters[i], true
		}
	}
	return Parameter{}, false
}

// ShouldExcludePath determines whether the given path should be removed (not rendered) from the template
func (c Config) ShouldExcludePath(path string) bool {
	if strings.HasSuffix(path, "t8.hcl") ||
		strings.HasSuffix(path, "t8.yml") ||
		strings.HasSuffix(path, "before.t8") ||
		strings.HasSuffix(path, "after.t8") {
		return true
	}

	for _, excludePath := range c.ExcludePaths {
		if isUnconditionalExcludePath(excludePath) {
			for _, excludePath := range excludePath.Paths {
				matched, _ := regexp.MatchString(excludePath, path)
				if matched {
					return true
				}
			}
		}

		for _, param := range c.Parameters {

			if isConditionalExcludePath(excludePath) {
				if excludePath.ParameterName == param.Name {

					if excludePath.Operator == NotEqual && excludePath.ParameterValue != param.Actual {
						for _, excludePath := range excludePath.Paths {
							matched, _ := regexp.MatchString(excludePath, path)
							if matched {
								return true
							}
						}
					}

					if excludePath.Operator == Equal && excludePath.ParameterValue == param.Actual {
						for _, excludePath := range excludePath.Paths {
							matched, _ := regexp.MatchString(excludePath, path)
							if matched {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

func isConditionalExcludePath(excludePath ExcludePath) bool {
	return excludePath.ParameterName != ""
}

func isUnconditionalExcludePath(excludePath ExcludePath) bool {
	return !isConditionalExcludePath(excludePath)
}

// DefaultString serializes the default parameter value as a string
func (p Parameter) DefaultString() string {
	if p.Default == nil {
		return ""
	}

	if p.Type == Option {
		t := p.Default.([]interface{})
		s := make([]string, len(t))
		for i, v := range t {
			s[i] = fmt.Sprint(v)
		}
		return strings.Join(s, ",")
	}
	return p.Default.(string)
}
