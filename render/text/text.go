// Package text provides tools to work with the textual representations of lambda terms.
package text

import (
	"fmt"

	"github.com/nightmarlin/lambda"
	"github.com/nightmarlin/lambda/render/text/internal"
)

func ValidateVariableName(v *lambda.Variable) error {
	if *v == "" {
		return &EmptyNameError{}
	}

	var err error
	for _, r := range *v {
		if !internal.ValidNameRune(r) {
			err = &InvalidNameError{Rune: r, Variable: v}
			break
		}
	}
	return err
}

type EmptyNameError struct{}

func (*EmptyNameError) Error() string {
	return "a variable name must contain at least one character"
}

type InvalidNameError struct {
	Rune     rune
	Variable *lambda.Variable
}

func (e *InvalidNameError) Error() string {
	return fmt.Sprintf("encountered restricted rune %q in variable name %q", string(e.Rune), string(*e.Variable))
}
