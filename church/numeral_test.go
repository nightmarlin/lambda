package church_test

import (
	"errors"
	"testing"

	. "github.com/nightmarlin/lambda"
	. "github.com/nightmarlin/lambda/church"
)

func TestNumeralConfig_Numeral(t *testing.T) {
	t.Parallel()

	var (
		f   = NewVariable(DefaultNumeralFuncVarName)
		arg = NewVariable(DefaultNumeralArgVarName)

		altFName   = "fn"
		altF       = NewVariable(altFName)
		altArgName = "arg"
		altArg     = NewVariable(altArgName)
		altCFG     = NumeralConfig{
			FuncVarName: altFName,
			ArgVarName:  altArgName,
		}
	)

	for name, entry := range map[string]struct {
		N    int
		CFG  NumeralConfig
		Want Term
	}{
		"zero": {
			N: 0,
			Want: NewAbstraction(
				f,
				NewAbstraction(
					arg,
					arg,
				),
			),
		},
		"zero, named": {
			N:   0,
			CFG: altCFG,
			Want: NewAbstraction(
				altF,
				NewAbstraction(
					altArg,
					altArg,
				),
			),
		},
		"one": {
			N: 1,
			Want: NewAbstraction(
				f,
				NewAbstraction(
					arg,
					NewApplication(f, arg),
				),
			),
		},
		"one, named": {
			N:   1,
			CFG: altCFG,
			Want: NewAbstraction(
				altF,
				NewAbstraction(
					altArg,
					NewApplication(altF, altArg),
				),
			),
		},
		"two": {
			N: 2,
			Want: NewAbstraction(
				f,
				NewAbstraction(
					arg,
					NewApplication(f, NewApplication(f, arg)),
				),
			),
		},
	} {
		t.Run(
			name,
			func(t *testing.T) {
				t.Parallel()

				got := entry.CFG.Numeral(entry.N)
				if !Equal(entry.Want, got) {
					t.Errorf(
						"%#v.Numeral(%d) = %v; want %v",
						entry.CFG, entry.N, got, entry.Want,
					)
				}
			},
		)
	}

	t.Run(
		"panics if n < 0",
		func(t *testing.T) {
			t.Parallel()

			var (
				n         = -2
				c         = make(chan struct{})
				recovered any
			)

			go func() {
				defer func() {
					recovered = recover()
					close(c)
				}()

				_ = Numeral(n)
			}()
			<-c

			if recovered == nil {
				t.Fatalf("expected Numeral(%d) to panic", n)
			}
			recoveredErr, ok := recovered.(error)
			if !ok {
				t.Fatalf("expected Numeral(%d) to panic with an error, got non-error value %#v", n, recovered)
			}

			ltzErr, ok := errors.AsType[*LessThanZeroError](recoveredErr)
			if !ok {
				t.Fatalf("expected Numeral(%d) to panic with a *LessThanZeroError, got %#v", n, recovered)
			}

			if ltzErr.N != n {
				t.Errorf("expected Numeral(%d) to panic with a *LessThanZeroError with value %d, got %d", n, n, ltzErr.N)
			}
		},
	)
}
