package lambda

import (
	"github.com/nightmarlin/lambda/internal/set"
)

// Closed returns whether every Variable used in the Term is bound by a parent Abstraction.
func Closed(term Term) bool { return closed(term, nil) }

func closed(term Term, bindings *set.Set[Variable]) bool {
	switch term := term.(type) {
	case *Variable:
		return bindings.Contains(*term)

	case *Abstraction:
		bindings, _ := bindings.Insert(*term.Parameter)
		return closed(term, bindings)

	case *Application:
		return closed(term.Receiver, bindings) &&
			closed(term.Argument, bindings)

	default:
		panic(&UnsupportedTermError{Term: term})
	}
}

// Cyclic returns whether a Term exists that refers to its own ancestor.
// Most functions will panic if they encounter a cycle, so this can be used to guard against that.
func Cyclic(term Term) bool { return cyclic(term, nil) }

func cyclic(term Term, seen *set.Set[Term]) bool {
	seen, isNew := seen.Insert(term)
	if !isNew {
		// cycle!
		return true
	}

	switch term := term.(type) {
	case *Variable:
		return false

	case *Abstraction:
		return cyclic(term, seen)

	case *Application:
		return cyclic(term.Receiver, seen) ||
			cyclic(term.Argument, seen)

	default:
		panic(&UnsupportedTermError{Term: term})
	}
}
