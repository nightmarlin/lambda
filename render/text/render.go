package text

import (
	"bytes"
	"fmt"
	"io"

	"github.com/nightmarlin/lambda"
	"github.com/nightmarlin/lambda/internal/set"
)

// Render writes the text form of the [lambda.Term] to the output stream.
//
//   - Variables are written directly.
//   - Abstractions are written as \<variable>.<body>
//   - Applications are written as (<receiver> <argument>).
//
// The output of Render may contain more parentheses than is strictly necessary.
// A trailing newline is included.
//
// [lambda.Variable] names must not contain any of the [RestrictedChars].
func Render(w io.Writer, term lambda.Term) error {
	if err := render(w, term, nil); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	return nil
}

func render(w io.Writer, term lambda.Term, seen *set.Set[lambda.Term]) error {
	seen, isNew := seen.Insert(term)
	if !isNew {
		return &lambda.CycleError{At: term}
	}

	type queueFn func() error

	var (
		q = func(fns ...queueFn) error {
			for _, f := range fns {
				if err := f(); err != nil {
					return err
				}
			}
			return nil
		}
		qPrint = func(s string) queueFn {
			return func() error {
				_, err := fmt.Fprint(w, s)
				return err
			}
		}
		qRender = func(e lambda.Term) queueFn {
			return func() error { return render(w, e, seen) }
		}
	)

	switch term := term.(type) {
	case *lambda.Variable:
		if err := ValidateVariableName(term); err != nil {
			return err
		}
		if err := q(qPrint(string(*term))); err != nil {
			return err
		}

	case *lambda.Abstraction:
		if err := q(
			qPrint("\\"),
			qRender(term.Parameter),
			qPrint("."),
			qRender(term.Body),
		); err != nil {
			return err
		}

	case *lambda.Application:
		if err := q(
			qPrint("("),
			qRender(term.Receiver),
			qPrint(" "),
			qRender(term.Argument),
			qPrint(")"),
		); err != nil {
			return err
		}

	default:
		panic(&lambda.UnsupportedTermError{Term: term})
	}

	return nil
}

// RenderString is equivalent to [Render] collected into a string,
// with the trailing newline removed.
func RenderString(term lambda.Term) (string, error) {
	var b bytes.Buffer
	if err := Render(&b, term); err != nil {
		return "", err
	}
	return string(b.Bytes()[:b.Len()-1]), nil
}
