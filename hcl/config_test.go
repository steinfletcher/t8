package hcl_test

import (
	"github.com/steinfletcher/t8/domain"
	"github.com/steinfletcher/t8/hcl"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestDecode_String(t *testing.T) {
	rawConfig := readFile("t8.hcl")

	conf, err := hcl.ParseConfig(rawConfig)

	assert.NoError(t, err)
	assert.Equal(t, "My Amazing App", conf.Name)
	assert.Equal(t, "ProjectName", conf.Parameters[0].Name)
	assert.Equal(t, "the project name", conf.Parameters[0].Description)
	assert.Equal(t, domain.String, conf.Parameters[0].Type)
	assert.Equal(t, "acme", conf.Parameters[0].Default)
}

func TestDecode_ExcludePaths(t *testing.T) {
	rawConfig := readFile("t8.hcl")

	conf, err := hcl.ParseConfig(rawConfig)

	assert.NoError(t, err)
	assert.Equal(t, "SqlDialect", conf.ExcludePaths[0].ParameterName)
	assert.Equal(t, domain.NotEqual, conf.ExcludePaths[0].Operator)
	assert.Equal(t, "postgresql", conf.ExcludePaths[0].ParameterValue)
	assert.Equal(t, []string{"^/postgres/.*$"}, conf.ExcludePaths[0].Paths)
}

func TestDecode_Option(t *testing.T) {
	rawConfig := readFile("t8.hcl")

	conf, err := hcl.ParseConfig(rawConfig)

	assert.NoError(t, err)
	assert.Equal(t, "SqlDialect", conf.Parameters[2].Name)
	assert.Equal(t, "the SQL dialect", conf.Parameters[2].Description)
	assert.Equal(t, domain.Option, conf.Parameters[2].Type)
	assert.Equal(t, []interface{}{"postgresql", "mysql"}, conf.Parameters[2].Default)
}

func TestDecode_YAML(t *testing.T) {
	rawConfig := readFile("t8.yml")

	conf, err := hcl.ParseConfig(rawConfig)

	assert.NoError(t, err)
	assert.Equal(t, "My Amazing App", conf.Name)
	assert.Equal(t, "ProjectName", conf.Parameters[0].Name)
	assert.Equal(t, "the project name", conf.Parameters[0].Description)
	assert.Equal(t, domain.String, conf.Parameters[0].Type)
	assert.Equal(t, "acme", conf.Parameters[0].Default)
}

func TestDecode_ErrorIfInvalidType(t *testing.T) {
	rawConfig := readFile("t8-invalid-type.hcl")

	_, err := hcl.ParseConfig(rawConfig)

	assert.EqualError(t, err, "'endofunctor' is not a valid value for type")
}

func TestDecode_ErrorIfInvalidHCL(t *testing.T) {
	rawConfig := readFile("t8-invalid-hcl.hcl")

	_, err := hcl.ParseConfig(rawConfig)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse HCL config")
}

func readFile(name string) []byte {
	data, err := ioutil.ReadFile("testdata/" + name)
	if err != nil {
		panic(err)
	}
	return data
}
