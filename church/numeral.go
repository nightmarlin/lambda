package church

import (
	"fmt"

	. "github.com/nightmarlin/lambda"
)

const (
	DefaultNumeralFuncVarName = "f"
	DefaultNumeralArgVarName  = "x"
)

type LessThanZeroError struct{ N int }

func (e *LessThanZeroError) Error() string {
	return fmt.Sprintf("church numerals cannot represent values less than 0: %d", e.N)
}

type NumeralConfig struct {
	FuncVarName string
	ArgVarName  string
}

// Numeral is equivalent to [NumeralConfig.Numeral](n) with no options set.
func Numeral(n int) *Abstraction { return NumeralConfig{}.Numeral(n) }

// Numeral creates the Church numeral n:
//
//	n == 0: \func.\arg.arg
//	n == 1: \func.\arg.(func arg)
//	n == 2: \func.\arg.(func (func arg))
//	and so on...
//
// https://en.wikipedia.org/wiki/Church_encoding#Church_numerals
// It panics if n < 0.
func (cfg NumeralConfig) Numeral(n int) *Abstraction {
	if n < 0 {
		panic(&LessThanZeroError{N: n})
	}
	if cfg.FuncVarName == "" {
		cfg.FuncVarName = DefaultNumeralFuncVarName
	}
	if cfg.ArgVarName == "" {
		cfg.ArgVarName = DefaultNumeralArgVarName
	}

	var (
		arg = &Abstraction{
			Parameter: NewVariable(cfg.ArgVarName),
			Body:      NewVariable(cfg.ArgVarName),
		}
		f = &Abstraction{
			Parameter: NewVariable(cfg.FuncVarName),
			Body:      arg,
		}
	)

	// wrap arg.Body with n applications to f
	for range n {
		arg.Body = &Application{
			Receiver: NewVariable(cfg.FuncVarName),
			Argument: arg.Body,
		}
	}

	return f
}
