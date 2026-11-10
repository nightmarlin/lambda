package lambda_test

import (
	"testing"

	. "github.com/nightmarlin/lambda"
)

// id returns the identity abstraction for the variable
func id(v *Variable) *Abstraction { return NewAbstraction(v, v) }

func TestBetaReduce(t *testing.T) {
	t.Parallel()

	var (
		a = NewVariable("a")
		b = NewVariable("b")
		x = NewVariable("x")
		y = NewVariable("y")
	)

	for name, entry := range map[string][]Term{
		// `a` cannot be reduced
		`variable identity`: {a},

		// `\a.a` cannot be reduced
		`abstraction identity`: {id(a)},

		`application`: {
			// (\x.x a)
			NewApplication(id(x), a),
			// a
			a,
		},

		// (a b) cannot be reduced
		`applications only be reduced when target is an abstraction`: {
			NewApplication(a, b),
		},

		`applications only modify abstractions, nested`: {
			// ((\x.x a) b)
			NewApplication(
				NewApplication(id(x), a),
				b,
			),
			// (a b)
			NewApplication(a, b),
		},

		`shadowed names`: {
			// (\x.(\x.x a) b)
			NewApplication(
				NewAbstraction(
					x,
					NewApplication(id(x), a),
				),
				b,
			),
			// (\x.x a)
			NewApplication(id(x), a),
			// a
			a,
		},

		`unbound applications`: {
			// (\x.\y.y a)
			NewApplication(
				NewAbstraction(x, id(y)),
				a,
			),
			// \y.y
			id(y),
		},
	} {
		t.Run(
			name,
			func(t *testing.T) {
				t.Parallel()

				if !Equal(entry[0], entry[0]) {
					t.Errorf("identity: (%s).Equal(%s) = false; want true", entry[0], entry[0])
				}

				var (
					reductions = 0
					current    = entry[0]
				)

				for {
					reductions += 1
					next, ok := BetaReduce(current)
					if !ok {
						break
					}

					if reductions >= len(entry) {
						t.Fatalf("still reducing after %d reductions", len(entry))
					}

					if want := entry[reductions]; !Equal(want, next) {
						t.Fatalf("reducing %s, got %s; want %s", current, next, want)
					}

					current = next
				}

				if reductions != len(entry) {
					t.Errorf("want %d reductions, performed %d", len(entry), reductions)
				}
			},
		)
	}
}
