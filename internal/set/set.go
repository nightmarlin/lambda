package set

import "maps"

// A Set is an immutable set implementation.
// Calls to Insert(OrPanic) construct a new Set with the old elements.
// A nil Set is safe for use.
type Set[T comparable] struct {
	s map[T]struct{}
}

func (s *Set[T]) Contains(v T) bool {
	if s == nil {
		return false
	}
	_, ok := s.s[v]
	return ok
}

func (s *Set[T]) Insert(v T) (next *Set[T], inserted bool) {
	oldLen := 0
	if s == nil || s.s == nil {
		next = &Set[T]{s: make(map[T]struct{})}
	} else {
		next = &Set[T]{s: maps.Clone(s.s)}
		oldLen = len(s.s)
	}
	next.s[v] = struct{}{}
	return next, len(next.s) != oldLen
}

func (s *Set[T]) InsertOrPanic(v T, orPanic error) *Set[T] {
	next, ok := s.Insert(v)
	if !ok {
		panic(orPanic)
	}
	return next
}
