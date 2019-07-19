package t8

//go:generate gonum -types=ParameterTypeEnum

type Validator func(value string) error

func ComposeValidators(validators ...Validator) Validator {
	return func(value string) error {
		for _, validate := range validators {
			if err := validate(value); err != nil {
				return err
			}
		}
		return nil
	}
}

type Config struct {
	Parameters []Parameter
	Name       string
	TargetDir  string
}

func (c *Config) GetParameter(name string) *Parameter {
	for _, v := range c.Parameters {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

type Parameter struct {
	Name        string
	Type        ParameterType
	Description string
	Default     interface{}
	Actual      interface{}
	Validators  []Validator
}

type ParameterTypeEnum struct {
	String string `enum:"string"`
	List   string `enum:"list"`
	Option string `enum:"option"`
}
