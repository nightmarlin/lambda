package text_test

import (
	"errors"
	"testing"

	. "github.com/nightmarlin/lambda"
	"github.com/nightmarlin/lambda/render/text"
)

func TestRender(t *testing.T) {
	t.Parallel()

	t.Run(
		"usage",
		func(t *testing.T) {
			t.Parallel()

			var (
				in = NewApplication(
					NewAbstraction(
						NewVariable("x"),
						NewVariable("x"),
					),
					NewVariable("a"),
				)

				want = "(\\x.x a)"

				got, err = text.RenderString(in)
			)
			if err != nil {
				t.Fatalf("must render: %v", err)
			}

			if got != want {
				t.Errorf("Render(%s) = %q; want %q", in, got, want)
			}
		},
	)

	t.Run(
		"cycle",
		func(t *testing.T) {
			t.Parallel()

			var (
				z    = NewVariable("z")
				abs3 = &Abstraction{Parameter: z, Body: z}

				y    = NewVariable("y")
				abs2 = &Abstraction{Parameter: y, Body: abs3}

				x    = NewVariable("x")
				abs1 = &Abstraction{Parameter: x, Body: abs2}
			)

			abs3.Body = abs1
			// \x.\y.\z.(_)
			//  ^--------'

			got, err := text.RenderString(abs1)
			if err == nil {
				t.Fatalf("must error, got %q", got)
			} else if got != "" {
				t.Errorf("failure should return no value")
			}

			asCycle, ok := errors.AsType[*CycleError](err)
			if !ok {
				t.Fatalf("error should be of type *CycleError, got %#v", err)
			} else if asCycle.At != abs1 {
				t.Errorf("cycle should be detected on second pass through abs1, got %#v", asCycle.At)
			}
		},
	)
}
