package domain

import (
	"errors"
)

var DefaultValidators = map[ParameterType][]Validator{
	String: {NotEmptyValidator},
	Option: {NewLengthValidator(1, 0)},
}

var NotEmptyValidator = func(s string) error {
	if s == "" {
		return errors.New("parameter may not be empty")
	}
	return nil
}

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

func NewLengthValidator(min, max int) Validator {
	return func(s string) error {
		if len(s) < min {
			return errors.New("length is less than 'min'")
		}

		if max != 0 && len(s) > max {
			return errors.New("length is greater than 'max'")
		}
		return nil
	}
}
