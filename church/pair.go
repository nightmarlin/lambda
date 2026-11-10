package church

import (
	. "github.com/nightmarlin/lambda"
)

var (
	DefaultPairFirstVarName    = "x"
	DefaultPairSecondVarName   = "y"
	DefaultPairSelectorVarName = "z"
)

type PairConfig struct {
	FirstVarName    string
	SecondVarName   string
	SelectorVarName string
}

// Pair is equivalent to [PairConfig.Pair]() with no options set.
func Pair() *Abstraction { return PairConfig{}.Pair() }

// Pair creates the Church pair constructor:
//
//	\first.\second.\selector.(selector (first second))
//
// The selector is compatible with [Boolean].
func (cfg PairConfig) Pair() *Abstraction {
	if cfg.FirstVarName == "" {
		cfg.FirstVarName = DefaultPairFirstVarName
	}
	if cfg.SecondVarName == "" {
		cfg.SecondVarName = DefaultPairSecondVarName
	}
	if cfg.SelectorVarName == "" {
		cfg.SelectorVarName = DefaultPairSelectorVarName
	}

	return NewAbstraction(
		NewVariable(cfg.FirstVarName),
		NewAbstraction(
			NewVariable(cfg.SecondVarName),
			NewAbstraction(
				NewVariable(cfg.SelectorVarName),
				NewApplication(
					NewVariable(cfg.SelectorVarName),
					NewApplication(NewVariable(cfg.FirstVarName), NewVariable(cfg.SecondVarName)),
				),
			),
		),
	)
}
