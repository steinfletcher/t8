package domain_test

import (
	"github.com/pkg/errors"
	"github.com/steinfletcher/t8/domain"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var config = domain.Config{
	Name: "my project",
	Parameters: []domain.Parameter{
		{
			Name:    "myName",
			Default: "myProject",
		},
	},
}

func TestConfig_GetParameter(t *testing.T) {
	cases := map[string]struct {
		parameterName     string
		expectedParameter domain.Parameter
		expectedBool      bool
	}{
		"success": {
			parameterName:     "myName",
			expectedParameter: domain.Parameter{Default: "myProject", Name: "myName"},
			expectedBool:      true,
		},
		"not found": {
			parameterName:     "doesn't Exist",
			expectedParameter: domain.Parameter{},
			expectedBool:      false,
		},
	}
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			parameter, ok := config.GetParameter(test.parameterName)

			assert.Equal(t, test.expectedParameter, parameter)
			assert.Equal(t, test.expectedBool, ok)
		})
	}
}

func TestComposeValidators(t *testing.T) {
	helloErr := errors.New("should contain hello")
	worldErr := errors.New("should contain world")

	helloValidator := func(in string) error {
		if strings.Contains(in, "hello") {
			return nil
		}
		return helloErr
	}

	worldValidator := func(in string) error {
		if strings.Contains(in, "world") {
			return nil
		}
		return worldErr
	}

	cases := map[string]struct {
		input       string
		expectedErr error
	}{
		"passes both validators": {
			input:       "hello world",
			expectedErr: nil,
		},
		"validator1 fail": {
			input:       "hello",
			expectedErr: worldErr,
		},
		"validator2 fail": {
			input:       "world",
			expectedErr: helloErr,
		},
	}
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			validate := domain.ComposeValidators(helloValidator, worldValidator)

			err := validate(test.input)

			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestParameter_String(t *testing.T) {
	cases := map[string]struct {
		param domain.Parameter
		out   string
	}{
		"option type with default": {
			param: domain.Parameter{
				Type:    domain.Option,
				Default: []interface{}{"a", "b"},
			},
			out: "a,b",
		},
		"without default": {
			param: domain.Parameter{
				Type: domain.Option,
			},
			out: "",
		},
		"string type": {
			param: domain.Parameter{
				Type:    domain.String,
				Default: "value",
			},
			out: "value",
		},
	}
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			s := test.param.DefaultString()
			assert.Equal(t, test.out, s)
		})
	}
}

func TestRemovePath_ShouldExcludePath(t *testing.T) {
	{
		cases := map[string]struct {
			path              string
			conf              domain.Config
			shouldExcludePath bool
		}{
			"should exclude unconditional path": {
				path: "test.sh",
				conf: domain.Config{
					ExcludePaths: []domain.ExcludePath{
						{
							Paths: []string{"abcdef", "^test.sh$"},
						},
					},
				},
				shouldExcludePath: true,
			},
			"should not exclude unconditional path": {
				path: "test.sh",
				conf: domain.Config{
					ExcludePaths: []domain.ExcludePath{
						{
							Paths: []string{"abcdef"},
						},
					},
				},
				shouldExcludePath: false,
			},
			"should exclude t8 hcl config file": {
				path:              "/blah/t8.hcl",
				shouldExcludePath: true,
			},
			"should exclude t8 yml config file": {
				path:              "/blah/t8.yml",
				shouldExcludePath: true,
			},
			"should exclude t8 before file": {
				path:              "before.t8",
				shouldExcludePath: true,
			},
			"should exclude t8 after file": {
				path:              "after.t8",
				shouldExcludePath: true,
			},
			"should exclude conditional path": {
				path: "/postgresql/query.go",
				conf: domain.Config{
					Parameters: []domain.Parameter{{Type: domain.String, Actual: "mysql", Name: "SqlDialect"}},
					ExcludePaths: []domain.ExcludePath{
						{
							ParameterName:  "SqlDialect",
							Operator:       domain.NotEqual,
							ParameterValue: "postgresql",
							Paths:          []string{"^/postgresql/.*$"},
						},
					},
				},
				shouldExcludePath: true,
			},
			"should exclude conditional path with multiple entries": {
				path: "/postgresql/query.go",
				conf: domain.Config{
					Parameters: []domain.Parameter{{Type: domain.String, Actual: "mysql", Name: "SqlDialect"}},
					ExcludePaths: []domain.ExcludePath{
						{
							ParameterName:  "SqlDialect",
							Operator:       domain.NotEqual,
							ParameterValue: "postgresql",
							Paths:          []string{"notmatching", "^/postgresql/.*.go$"},
						},
					},
				},
				shouldExcludePath: true,
			},
			"should not exclude path if no exclude path param": {
				path: "/somepath",
				conf: domain.Config{
					Parameters: []domain.Parameter{{Type: domain.String, Actual: "mysql", Name: "SqlDialect"}},
				},
				shouldExcludePath: false,
			},
			"should not exclude path if value correct but path does not match": {
				path: "/somepath",
				conf: domain.Config{
					Parameters: []domain.Parameter{{Type: domain.String, Actual: "mysql", Name: "SqlDialect"}},
					ExcludePaths: []domain.ExcludePath{
						{
							ParameterName:  "SqlDialect",
							Operator:       domain.NotEqual,
							ParameterValue: "postgresql",
							Paths:          []string{"^/postgres/.*$"},
						},
					},
				},
				shouldExcludePath: false,
			},
		}
		for name, test := range cases {
			t.Run(name, func(t *testing.T) {
				shouldRender := test.conf.ShouldExcludePath(test.path)

				assert.Equal(t, test.shouldExcludePath, shouldRender)
			})
		}
	}
}
