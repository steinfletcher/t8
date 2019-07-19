package config

import (
	"github.com/hashicorp/hcl"
	"github.com/pkg/errors"
	"github.com/steinfletcher/t8"
	"github.com/steinfletcher/t8/validation"
	"gopkg.in/yaml.v2"
)

type config struct {
	Name      string
	Parameter []hclParameter `hcl:"parameter" yaml:"parameter"`
}

type hclParameter map[string][]hclParameterValue

type hclParameterValue struct {
	Type        string      `hcl:"type" yaml:"type"`
	Description string      `hcl:"description" yaml:"description"`
	Default     interface{} `hcl:"default" yaml:"default"`
}

func New(rawConfig []byte) (t8.Config, error) {
	var p config
	hclErr := hcl.Decode(&p, string(rawConfig))
	if hclErr != nil {
		err := yaml.Unmarshal(rawConfig, &p)
		if err != nil {
			return t8.Config{}, errors.Wrap(hclErr, "failed to parse HCL config")
		}
	}
	return mapParams(p)
}

func mapParams(hclConfig config) (t8.Config, error) {
	var p []t8.Parameter

	for _, params := range hclConfig.Parameter {
		for k, v := range params {
			if len(v) > 0 {
				firstParam := v[0]
				typ, err := t8.NewParameterType(firstParam.Type)
				if err != nil {
					return t8.Config{}, err
				}

				p = append(p, t8.Parameter{
					Name:        k,
					Description: firstParam.Description,
					Type:        typ,
					Default:     firstParam.Default,
					Validators:  validation.DefaultValidators[typ],
				})
			}
		}
	}

	return t8.Config{
		Parameters: p,
		Name:       hclConfig.Name,
		TargetDir:  "repoz",
	}, nil
}
