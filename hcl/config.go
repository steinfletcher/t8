package hcl

import (
	"github.com/hashicorp/hcl"
	"github.com/pkg/errors"
	"github.com/steinfletcher/t8/domain"
	"gopkg.in/yaml.v2"
	"log"
)

type config struct {
	Name       string
	Parameter  []hclParameter  `hcl:"parameter"  yaml:"parameter"`
	RemovePath []hclRemovePath `hcl:"excludePath" yaml:"excludePath"`
}

type hclParameter map[string][]hclParameterValue

type hclParameterValue struct {
	Type        string      `hcl:"type" yaml:"type"`
	Description string      `hcl:"description" yaml:"description"`
	Default     interface{} `hcl:"default" yaml:"default"`
}

type hclRemovePath map[string][]hclRemovePathValue

type hclRemovePathValue struct {
	MatchPath      []string `hcl:"paths" yaml:"paths"`
	ParameterValue string   `hcl:"parameterValue" yaml:"parameterValue"`
	ParameterName  string   `hcl:"parameterName" yaml:"parameterName"`
	Operator       string   `hcl:"operator" yaml:"operator"`
}

func ParseConfig(rawConfig []byte) (domain.Config, error) {
	var p config
	hclErr := hcl.Decode(&p, string(rawConfig))
	if hclErr != nil {
		err := yaml.Unmarshal(rawConfig, &p)
		if err != nil {
			return domain.Config{}, errors.Wrap(hclErr, "failed to parse HCL config")
		}
	}
	return mapParams(p)
}

func mapParams(hclConfig config) (domain.Config, error) {
	var params []domain.Parameter
	var excludePaths []domain.ExcludePath

	for _, p := range hclConfig.Parameter {
		for k, v := range p {
			if len(v) > 0 {
				firstParam := v[0]
				typ, err := domain.NewParameterType(firstParam.Type)
				if err != nil {
					return domain.Config{}, err
				}

				params = append(params, domain.Parameter{
					Name:        k,
					Description: firstParam.Description,
					Type:        typ,
					Default:     firstParam.Default,
					Validators:  domain.DefaultValidators[typ],
				})
			}
		}
	}

	for _, paths := range hclConfig.RemovePath {
		for _, v := range paths {
			if len(v) > 0 {
				firstParam := v[0]
				path := domain.ExcludePath{
					ParameterName:  firstParam.ParameterName,
					ParameterValue: firstParam.ParameterValue,
					Paths:          firstParam.MatchPath,
				}

				var err error
				if firstParam.Operator != "" {
					path.Operator, err = domain.NewRemovePathOperator(firstParam.Operator)
					if err != nil {
						return domain.Config{}, err
					}
				}

				excludePaths = append(excludePaths, path)
			}
		}
	}

	conf := domain.Config{
		Parameters:   params,
		ExcludePaths: excludePaths,
		Name:         hclConfig.Name,
	}

	if domain.IsDebug {
		log.Printf("Parsed config: %+v\n", conf)
	}

	return conf, nil
}
