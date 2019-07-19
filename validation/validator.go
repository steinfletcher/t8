package validation

import (
	"errors"
	"github.com/steinfletcher/t8"
)

var DefaultValidators = map[t8.ParameterType][]t8.Validator{
	t8.String: {NotEmptyValidator},
	t8.List:   {NewLengthValidator(1, 0)},
}

var NotEmptyValidator = func(s string) error {
	if s == "" {
		return errors.New("parameter may not be empty")
	}
	return nil
}

func NewLengthValidator(min, max int) t8.Validator {
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
