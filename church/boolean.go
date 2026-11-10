package church

import (
	. "github.com/nightmarlin/lambda"
)

const (
	DefaultBooleanTruthyVarName = "t"
	DefaultBooleanFalsyVarName  = "f"
)

type BooleanConfig struct {
	TruthyVarName string
	FalsyVarName  string
}

// Boolean is equivalent to [BooleanConfig.Boolean](b) with no options set.
func Boolean(b bool) *Abstraction { return BooleanConfig{}.Boolean(b) }

// True is equivalent to [Boolean](true).
// It is also equivalent to the Church pair selector `first`.
func True() *Abstraction { return Boolean(true) }

// False is equivalent to [Boolean](false).
// It is also equivalent to the Church pair selector `second`.
func False() *Abstraction { return Boolean(false) }

// Boolean creates the Church boolean (and choice) b:
//
//	true: \truthy.\falsy.truthy
//	false: \truthy.\falsy.falsy
//
// Boolean can be used as a Church pair selector with [True] as `first` and [False] as `second`.
func (cfg BooleanConfig) Boolean(b bool) *Abstraction {
	if cfg.TruthyVarName == "" {
		cfg.TruthyVarName = DefaultBooleanTruthyVarName
	}
	if cfg.FalsyVarName == "" {
		cfg.FalsyVarName = DefaultBooleanFalsyVarName
	}

	out := cfg.TruthyVarName
	if !b {
		out = cfg.FalsyVarName
	}

	// \t.\f.(b ? t : f)
	return NewAbstraction(
		NewVariable(cfg.TruthyVarName),
		NewAbstraction(
			NewVariable(cfg.FalsyVarName),
			NewVariable(out),
		),
	)
}

// True is equivalent to [BooleanConfig.Boolean](true).
// It is also equivalent to the Church pair selector `first`.
func (cfg BooleanConfig) True() *Abstraction { return cfg.Boolean(true) }

// False is equivalent to [BooleanConfig.Boolean](false).
// It is also equivalent to the Church pair selector `second`.
func (cfg BooleanConfig) False() *Abstraction { return cfg.Boolean(false) }
