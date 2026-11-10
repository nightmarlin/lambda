package internal

import (
	"errors"
	"io"

	"github.com/nightmarlin/lambda"
)

func Parse(r io.Reader) (lambda.Term, error) {
	var (
		lex   = lNewLexer(r)
		parse = lNewParser()
	)
	if parse.Parse(lex) != 0 {
		return nil, errors.Join(lex.errors...)
	}
	return lex.expression, nil
}

// Verbosity sets the global parser log level.
// Calling `reset` returns it to the value before Verbosity was called.
// `reset` is compatible with [testing.T.Cleanup].
func Verbosity(lv int) (reset func()) {
	old := lDebug
	lDebug = lv
	return func() { lDebug = old }
}
