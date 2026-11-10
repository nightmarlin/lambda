package church

import (
	. "github.com/nightmarlin/lambda"
)

const (
	DefaultOptionalFailureVarName = "f"
	DefaultOptionalSuccessVarName = "s"
)

type OptionalConfig struct {
	SuccessVarName string
	FailureVarName string
}

// None is equivalent to [OptionalConfig.None]() with no options set.
func None() *Abstraction { return OptionalConfig{}.None() }

// Some is equivalent to [OptionalConfig.Some](v) with no options set.
func Some(v *Variable) *Abstraction { return OptionalConfig{}.Some(v) }

// None creates the "empty" Church optional:
//
//	\failure.\success.failure
func (cfg OptionalConfig) None() *Abstraction {
	if cfg.SuccessVarName == "" {
		cfg.SuccessVarName = DefaultOptionalSuccessVarName
	}
	if cfg.FailureVarName == "" {
		cfg.FailureVarName = DefaultOptionalFailureVarName
	}

	return NewAbstraction(
		NewVariable(cfg.FailureVarName),
		NewAbstraction(
			NewVariable(cfg.SuccessVarName),
			NewVariable(cfg.FailureVarName),
		),
	)
}

// Some creates the "filled" Church optional:
//
//	\failure.\success.(success v)
//
// v is assumed to be bound by some parent [lambda.Abstraction].
func (cfg OptionalConfig) Some(v *Variable) *Abstraction {
	if cfg.SuccessVarName == "" {
		cfg.SuccessVarName = DefaultOptionalSuccessVarName
	}
	if cfg.FailureVarName == "" {
		cfg.FailureVarName = DefaultOptionalFailureVarName
	}

	return NewAbstraction(
		NewVariable(cfg.FailureVarName),
		NewAbstraction(
			NewVariable(cfg.SuccessVarName),
			NewApplication(NewVariable(cfg.SuccessVarName), v),
		),
	)
}
