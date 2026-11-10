package lambda_test

import (
	"testing"

	. "github.com/nightmarlin/lambda"
)

func TestAlphaEqual(t *testing.T) {
	t.Parallel()

	var (
		a = NewVariable("a")
		b = NewVariable("b")
		c = NewVariable("c")
		x = NewVariable("x")
		y = NewVariable("y")
	)

	for name, entry := range map[string]struct {
		A, B Term
		Want bool
	}{
		"free vars are the same":     {A: a, B: a, Want: true},
		"free vars must be the same": {A: a, B: b, Want: false},
		"id lambda is the same":      {A: id(a), B: id(b), Want: true},

		"fully bound tree is the same": {
			A:    NewAbstraction(a, NewApplication(NewAbstraction(x, id(x)), a)),
			B:    NewAbstraction(b, NewApplication(NewAbstraction(y, id(y)), b)),
			Want: true,
		},

		"tree with free var is the same": {
			A:    NewAbstraction(a, NewApplication(NewAbstraction(x, id(x)), c)),
			B:    NewAbstraction(b, NewApplication(NewAbstraction(y, id(y)), c)),
			Want: true,
		},

		"tree without shadowing is not the same": {
			// if an arbitrary value passed to `a` uses `x`, the meaning changes
			A:    NewAbstraction(x, NewApplication(x, NewAbstraction(y, NewAbstraction(x, NewAbstraction(x, y))))),
			B:    NewAbstraction(x, NewApplication(x, NewAbstraction(y, NewAbstraction(a, NewAbstraction(a, y))))),
			Want: false,
		},
		"tree with shadowing is not the same": {
			A:    NewAbstraction(x, NewApplication(x, NewAbstraction(y, NewAbstraction(x, NewAbstraction(x, y))))),
			B:    NewAbstraction(a, NewApplication(a, NewAbstraction(y, NewAbstraction(a, NewAbstraction(a, y))))),
			Want: true,
		},
	} {
		t.Run(
			name,
			func(t *testing.T) {
				t.Parallel()

				got := AlphaEqual(entry.A, entry.B)
				if got != entry.Want {
					t.Errorf(
						"AlphaEqual(%s, %s) = %t; want %t",
						entry.A, entry.B, got, entry.Want,
					)
				}
			},
		)
	}
}
