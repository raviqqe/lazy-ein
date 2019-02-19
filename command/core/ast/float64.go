package ast

import "github.com/ein-lang/ein/command/core/types"

// Float64 is a float64 literal.
type Float64 struct {
	value float64
}

// NewFloat64 creates a float64 number.
func NewFloat64(f float64) Float64 {
	return Float64{f}
}

// Value returns a value.
func (f Float64) Value() float64 {
	return f.value
}

// ConvertTypes converts types.
func (f Float64) ConvertTypes(func(types.Type) types.Type) Expression {
	return f
}

// RenameVariables renames variables.
func (f Float64) RenameVariables(map[string]string) Expression {
	return f
}

// RenameVariablesInAtom renames variables.
func (f Float64) RenameVariablesInAtom(map[string]string) Atom {
	return f
}

func (Float64) isAtom()       {}
func (Float64) isExpression() {}
