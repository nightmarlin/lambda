package text

import (
	"io"
	"strings"

	"github.com/nightmarlin/lambda"
	"github.com/nightmarlin/lambda/render/text/internal"
)

// Parse extracts a [lambda.Term] from the input stream.
//
// Terms are parsed as follows:
//
//	term        ::= variable | abstraction | application | '(' term ')'
//	variable    ::= (UTF-8 letters, numbers, '_' & '-')
//	abstraction ::= '\' variable '.' term
//	application ::= '(' term term ')'
func Parse(r io.Reader) (lambda.Term, error) { return internal.Parse(r) }

// ParseString is equivalent to Parse(strings.NewReader(s)).
func ParseString(s string) (lambda.Term, error) { return Parse(strings.NewReader(s)) }
