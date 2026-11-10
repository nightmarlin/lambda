package lambda

import (
	"maps"
)

// Equal checks that a and b are identical:
// All structures match and all variable names are the same.
func Equal(a, b Term) bool {
	switch a := a.(type) {
	case *Variable:
		bVar, ok := b.(*Variable)
		return ok && *a == *bVar

	case *Abstraction:
		bAbs, ok := b.(*Abstraction)
		return ok &&
			Equal(a.Parameter, bAbs.Parameter) &&
			Equal(a.Body, bAbs.Body)

	case *Application:
		bApp, ok := b.(*Application)
		return ok &&
			Equal(a.Argument, bApp.Argument) &&
			Equal(a.Receiver, bApp.Receiver)

	default:
		panic(&UnsupportedTermError{Term: a})
	}
}

// AlphaEqual checks that a and b are structurally equal:
// If bound variables in a and b were renamed arbitrarily (preserving shadows),
// the two terms would still have the same behaviour.
func AlphaEqual(a, b Term) bool {
	return alphaEqual(a, b, newAlphaScopes())
}

func alphaEqual(a, b Term, scopes *alphaScopes) bool {
	switch a := a.(type) {
	// if a is bound, b must be bound to the equivalent scope
	// if a is free, b must also be free, and they must refer to the same free variable
	case *Variable:
		b, ok := b.(*Variable)
		if !ok {
			return false
		}
		return scopes.same(a, b)

	// if a binds a free var, b must also bind a free var.
	// if a shadows a bound var, b must shadow its equivalent var.
	// the bodies must also be equivalent.
	case *Abstraction:
		b, ok := b.(*Abstraction)
		if !ok {
			return false
		}
		s, ok := scopes.with(a.Parameter, b.Parameter)
		if !ok {
			return false
		}
		return alphaEqual(a.Body, b.Body, s)

	// the left and right sides must be equivalent.
	case *Application:
		b, ok := b.(*Application)
		if !ok {
			return false
		}

		return alphaEqual(a.Receiver, b.Receiver, scopes) &&
			alphaEqual(a.Argument, b.Argument, scopes)

	default:
		panic(&UnsupportedTermError{Term: a})
	}
}

func AlphaRename(from, to Term) (Term, bool) {
	// if two terms are alpha-equivalent, renaming from->to is the same as cloning to.
	if !AlphaEqual(from, to) {
		return nil, false
	}
	return to.Clone(), true
}

// region scope tracking

// alphaScopes provides variable scope tracking for alpha-equivalence.
type alphaScopes struct{ aToB, bToA map[Variable]Variable }

func newAlphaScopes() *alphaScopes {
	return &alphaScopes{
		aToB: make(map[Variable]Variable),
		bToA: make(map[Variable]Variable),
	}
}

// match looks up the bindings of a and b in s, and returns their equivalence:
//
//	(true, true):  a and b are both bound to the same scope
//	(false, true): a and b are both unbound
//	(?, false):    only one of a or b is bound, or both are bound but to different scopes
func (s *alphaScopes) match(a, b *Variable) (scope, broadly bool) {
	var (
		bInA, inA = s.aToB[*a]
		aInB, inB = s.bToA[*b]
	)

	switch {
	// unbound on both sides
	case !inA && !inB:
		return false, true

	// both are bound, and the equivalent scope is as expected
	case inA && inB &&
		Equal(&aInB, a) &&
		Equal(&bInA, b):
		return true, true

	// only one of a or b is bound, or
	// both are bound but to different scopes
	default:
		return false, false
	}
}

// same returns true when a and b are bound to the same scope,
// or unbound but refer to the same variable.
func (s *alphaScopes) same(a, b *Variable) bool {
	bound, both := s.match(a, b)
	if !both {
		return false
	}
	return bound || Equal(a, b)
}

// with attempts to add a and b to the scope, returning a new set on success.
// It fails if
//   - only one of a or b is already bound, or
//   - a and b are already bound to different scopes.
func (s *alphaScopes) with(a, b *Variable) (*alphaScopes, bool) {
	if _, ok := s.match(a, b); !ok {
		return nil, false
	}

	next := &alphaScopes{
		aToB: maps.Clone(s.aToB),
		bToA: maps.Clone(s.bToA),
	}

	next.aToB[*a] = *b
	next.bToA[*b] = *a

	return next, true
}

// endregion
