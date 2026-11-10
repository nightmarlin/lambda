// Package lambda implements a basic form of the lambda calculus.
package lambda

import (
	"fmt"

	"github.com/nightmarlin/lambda/internal/set"
)

// A Term is an element of a Lambda expression.
// Unless otherwise specified, functions that accept a Term (including methods in Term)
// are not expected to protect against cycles (references to a Term's own ancestor).
type Term interface {
	// All Term impls are pointers so we only need to check for *T values.

	term() // sealed

	// Clone creates a deep copy of the Term.
	// As Terms are passed by reference,
	// this should be used to avoid accidental mutations of the tree.
	// Clone always returns the same concrete type.
	Clone() Term
}

// A CycleError may be returned as an error or a panic in operations that do not support it.
type CycleError struct{ At Term }

func (e *CycleError) Error() string {
	return fmt.Sprintf("cycle found starting at %#v", e.At)
}

// An UnsupportedTermError is always returned as a panic
// if an unsupported term type is encountered.
// This should never happen.
type UnsupportedTermError struct{ Term Term }

func (e *UnsupportedTermError) Error() string {
	return fmt.Sprintf("unsupported term type: %#v", e.Term)
}

// A Variable is the basic unit of a lambda expression.
// It represents some value.
// A Variable may be considered
//
//   - Bound: (an [Abstraction] exists higher in the [Term]
//     that has the variable as its parameter), or
//
//   - Free: No binding [Abstraction] exists in the parent [Term].
type Variable string

func NewVariable[V ~string](v V) *Variable { return new(Variable(v)) }

func (v *Variable) term()       {}
func (v *Variable) Clone() Term { return NewVariable(*v) }

// An Abstraction binds the Parameter in the context of the Body.
// If the Parameter has the same name as another [Variable] in the parent [Term],
// this Parameter now takes precedence (shadows it).
type Abstraction struct {
	Parameter *Variable
	Body      Term
}

// NewAbstraction creates an abstraction binding the param in the body.
func NewAbstraction(param *Variable, body Term) *Abstraction {
	return &Abstraction{Parameter: param.Clone().(*Variable), Body: body.Clone()}
}

func (abs *Abstraction) term()       {}
func (abs *Abstraction) Clone() Term { return NewAbstraction(abs.Parameter, abs.Body) }

// Bind replaces all unshadowed instances of Parameter in Body with arg,
// returning the updated Body.
// It panics if it encounters a cycle.
func (abs *Abstraction) Bind(arg Term) Term {
	return bind(abs.Parameter, arg, abs.Body, nil)
}

func bind(find *Variable, replace Term, in Term, seen *set.Set[Term]) Term {
	seen = seen.InsertOrPanic(in, &CycleError{At: in})

	switch in := in.(type) {
	case *Variable:
		if !Equal(in, find) {
			return in.Clone()
		}
		return replace.Clone()

	case *Abstraction:
		// shadowed
		if Equal(in.Parameter, find) {
			return in.Clone()
		}
		return NewAbstraction(in.Parameter, bind(find, replace, in.Body, seen))

	case *Application:
		return NewApplication(bind(find, replace, in.Receiver, seen), bind(find, replace, in.Argument, seen))

	default:
		panic(&UnsupportedTermError{Term: in})
	}
}

// An Application describes the binding of the Argument to the Receiver,
// though it only has an effect if Receiver is an [Abstraction].
type Application struct{ Receiver, Argument Term }

func NewApplication(target, arg Term) *Application {
	return &Application{Receiver: target.Clone(), Argument: arg.Clone()}
}

func (app *Application) term()       {}
func (app *Application) Clone() Term { return NewApplication(app.Receiver, app.Argument) }
func (app *Application) Reducable() bool {
	_, ok := app.Receiver.(*Abstraction)
	return ok
}
