package text_test

import (
	"testing"

	"github.com/nightmarlin/lambda"
	"github.com/nightmarlin/lambda/render/text"
)

func TestParse(t *testing.T) {
	t.Parallel()

	// if uncommented, also uncomment t.Parallel(): verbosity is a global setting
	// t.Cleanup(internal.Verbosity(4))

	var (
		a = lambda.NewVariable("a")
		b = lambda.NewVariable("b")
		x = lambda.NewVariable("x")
		y = lambda.NewVariable("y")
	)

	for name, entry := range map[string]struct {
		Expr string
		Want lambda.Term
	}{
		`variable`:                        {Expr: `a`, Want: a},
		`abstraction`:                     {Expr: `\a.a`, Want: lambda.NewAbstraction(a, a)},
		`application`:                     {Expr: `(a a)`, Want: lambda.NewApplication(a, a)},
		`application without inner paren`: {Expr: `(a a)`, Want: lambda.NewApplication(a, a)},
		`kitchen sink`:                    {Expr: `(\x.x a)`, Want: lambda.NewApplication(lambda.NewAbstraction(x, x), a)},
		`nested apply`: {
			Expr: `((\x.\y.(x y) a) b)`,
			Want: lambda.NewApplication(
				lambda.NewApplication(
					lambda.NewAbstraction(
						x,
						lambda.NewAbstraction(
							y,
							lambda.NewApplication(x, y),
						),
					),
					a,
				),
				b,
			),
		},
	} {
		t.Run(
			name,
			func(t *testing.T) {
				t.Parallel()

				got, err := text.ParseString(entry.Expr)
				if err != nil {
					t.Fatalf("must parse: %s", err.Error())
				}

				if !lambda.Equal(entry.Want, got) {
					t.Errorf(
						"Parse(%q) = %s; want %s",
						entry.Expr, got, entry.Want,
					)
				}
			},
		)
	}
}
