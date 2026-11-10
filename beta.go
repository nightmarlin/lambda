package lambda

import (
	"fmt"
)

// BetaReduce performs a single Beta-Reduction operation,
// returning the resulting Term and whether any changes were made
// (IE when an *Application with a receiver of type *Abstraction was present in the tree).
// When multiple reducable *Application terms are present, any one may be reduced.
// https://en.wikipedia.org/wiki/Lambda_calculus_definition#%CE%B2-reduction
func BetaReduce(term Term) (Term, bool) {
	switch term := term.(type) {
	// variables can't be reduced
	case *Variable:
		break

	// try to reduce the body
	case *Abstraction:
		reducedBody, ok := BetaReduce(term.Body)
		if ok {
			return NewAbstraction(term.Parameter, reducedBody), true
		}

	// applications can be reduced directly if their receiver is an abstraction
	// by binding the argument to the receiver.
	// otherwise, try to reduce the receiver or argument.
	case *Application:
		if term.Reducable() {
			return term.Receiver.(*Abstraction).Bind(term.Argument), true
		}

		reducedReceiver, ok := BetaReduce(term.Receiver)
		if ok {
			return NewApplication(reducedReceiver, term.Argument), true
		}

		reducedArgument, ok := BetaReduce(term.Argument)
		if ok {
			return NewApplication(term.Receiver, reducedArgument), true
		}

	default:
		panic(fmt.Errorf("unsupported term type: %#v", term))
	}

	return nil, false
}
