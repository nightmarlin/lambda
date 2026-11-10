package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode"

	"github.com/nightmarlin/lambda"
)

//go:generate go tool golang.org/x/tools/cmd/goyacc -l -p l -o parser.gen.go -v parser.gen.output lambda.y

// always put details in error messages
func init() { lErrorVerbose = true }

var NameUnicodeRanges = []*unicode.RangeTable{unicode.Letter}

func ValidNameRune(r rune) bool {
	switch r {
	case '_', '-':
		return true
	default:
		return unicode.In(r, unicode.Letter, unicode.Number)
	}
}

func tokenFromRune(r rune) (int, bool) {
	switch r {
	case '(':
		return T_LEFT, true
	case ')':
		return T_RIGHT, true
	case '.':
		return T_DOT, true
	case '\\':
		return T_LAMBDA, true
	case ',':
		return T_COMMA, true
	default:
		if !ValidNameRune(r) {
			return 0, false
		}
		return T_NAME, true
	}
}

type lLexerImpl struct {
	r *bufio.Reader

	expression lambda.Term // result
	errors     []error
}

func lNewLexer(r io.Reader) *lLexerImpl {
	return &lLexerImpl{r: bufio.NewReader(r)}
}

func (l *lLexerImpl) Error(s string) {
	l.errors = append(l.errors, fmt.Errorf("%s", s))
}

// read consumes the next rune from the input stream.
func (l *lLexerImpl) read() (rune, error) {
	r, _, err := l.r.ReadRune()
	if err != nil {
		return 0, fmt.Errorf("read: %w", err)
	}
	return r, nil
}

// peek pulls the next rune from the input stream, then unreads it for later reuse.
func (l *lLexerImpl) peek() (rune, error) {
	r, _, err := l.r.ReadRune()
	if err != nil {
		return 0, fmt.Errorf("peek: read: %w", err)
	}

	if err := l.r.UnreadRune(); err != nil {
		return 0, fmt.Errorf("peek: unread: %w", err)
	}
	return r, nil
}

func (l *lLexerImpl) Lex(lval *lSymType) (res int) {
	var ongoing string // for T_NAME

	for {
		r, err := l.read()
		switch {
		case errors.Is(err, io.EOF):
			return 0
		case err != nil:
			l.errors = append(l.errors, err)
			return lErrCode
		}

		tok, ok := tokenFromRune(r)
		switch {
		// non-tokenizable chars are ignored
		case !ok:
			continue

		// text tokens may run for multiple runes
		case tok == T_NAME:
			ongoing = fmt.Sprintf("%s%s", ongoing, string(r))

			peeked, err := l.peek()
			if errors.Is(err, io.EOF) {
				lval.Name = ongoing
				return T_NAME
			} else if err != nil {
				l.errors = append(l.errors, err)
				return lErrCode
			}

			peekedToken, ok := tokenFromRune(peeked)
			if !ok || peekedToken != T_NAME {
				// if upcoming rune can't be part of T_TEXT, end here
				lval.Name = ongoing
				return T_NAME
			}

		// all other tokens are single-rune
		default:
			return tok
		}
	}
}
