package config_test

import (
	"github.com/steinfletcher/t8"
	"github.com/steinfletcher/t8/config"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestDecode_String(t *testing.T) {
	rawConfig := readFile("t8.hcl")

	conf, err := config.New(rawConfig)

	assert.NoError(t, err)
	assert.Equal(t, "My Amazing App", conf.Name)
	assert.Equal(t, "project_name", conf.Parameters[0].Name)
	assert.Equal(t, "the project name", conf.Parameters[0].Description)
	assert.Equal(t, t8.String, conf.Parameters[0].Type)
	assert.Equal(t, "acme", conf.Parameters[0].Default)
}

func TestDecode_List(t *testing.T) {
	rawConfig := readFile("t8.hcl")

	conf, err := config.New(rawConfig)

	assert.NoError(t, err)
	assert.Equal(t, "authors", conf.Parameters[1].Name)
	assert.Equal(t, "the project authors", conf.Parameters[1].Description)
	assert.Equal(t, t8.List, conf.Parameters[1].Type)
	assert.Equal(t, []interface{}{"yuki@foo.com", "mei@bar.com"}, conf.Parameters[1].Default)
}

func TestDecode_Option(t *testing.T) {
	rawConfig := readFile("t8.hcl")

	conf, err := config.New(rawConfig)

	assert.NoError(t, err)
	assert.Equal(t, "sql_dialect", conf.Parameters[2].Name)
	assert.Equal(t, "the SQL dialect", conf.Parameters[2].Description)
	assert.Equal(t, t8.Option, conf.Parameters[2].Type)
	assert.Equal(t, []interface{}{"postgresql", "mysql"}, conf.Parameters[2].Default)
}

func TestDecode_YAML(t *testing.T) {
	rawConfig := readFile("t8.yml")

	conf, err := config.New(rawConfig)

	assert.NoError(t, err)
	assert.Equal(t, "My Amazing App", conf.Name)
	assert.Equal(t, "project_name", conf.Parameters[0].Name)
	assert.Equal(t, "the project name", conf.Parameters[0].Description)
	assert.Equal(t, t8.String, conf.Parameters[0].Type)
	assert.Equal(t, "acme", conf.Parameters[0].Default)
}

func TestDecode_ErrorIfInvalidType(t *testing.T) {
	rawConfig := readFile("t8-invalid-type.hcl")

	_, err := config.New(rawConfig)

	assert.EqualError(t, err, "'endofunctor' is not a valid value for type")
}

func TestDecode_ErrorIfInvalidHCL(t *testing.T) {
	rawConfig := readFile("t8-invalid-hcl.hcl")

	_, err := config.New(rawConfig)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse HCL config")
}

func TestConfig_GetParameter(t *testing.T) {
	rawConfig := readFile("t8.hcl")
	conf, _ := config.New(rawConfig)

	parameter := conf.GetParameter("authors")

	assert.NotNil(t, parameter)
}

func readFile(name string) []byte {
	data, err := ioutil.ReadFile("testdata/" + name)
	if err != nil {
		panic(err)
	}
	return data
}
