package domain_test

import (
	"errors"
	"github.com/steinfletcher/t8/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNotEmptyValidator(t *testing.T) {
	cases := map[string]struct {
		input       string
		expectedErr error
	}{
		"valid":     {input: "content", expectedErr: nil},
		"not valid": {input: "", expectedErr: errors.New("parameter may not be empty")},
	}
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := domain.NotEmptyValidator(test.input)

			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestLengthValidator(t *testing.T) {
	cases := map[string]struct {
		input       string
		min         int
		max         int
		expectedErr error
	}{
		"valid":           {min: 1, max: 5, input: "ab", expectedErr: nil},
		"valid max limit": {min: 3, max: 5, input: "abcde", expectedErr: nil},
		"valid min limit": {min: 1, max: 3, input: "a", expectedErr: nil},
		"invalid min":     {min: 2, max: 3, input: "a", expectedErr: errors.New("length is less than 'min'")},
		"invalid max":     {min: 1, max: 3, input: "abcd", expectedErr: errors.New("length is greater than 'max'")},
	}
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			validator := domain.NewLengthValidator(test.min, test.max)

			err := validator(test.input)

			assert.Equal(t, test.expectedErr, err)
		})
	}
}
