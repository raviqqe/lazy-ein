package ast

import "github.com/ein-lang/ein/command/core/types"

// PrimitiveAlternative is a primitive alternative.
type PrimitiveAlternative struct {
	literal    Literal
	expression Expression
}

// NewPrimitiveAlternative creates a primitive alternative.
func NewPrimitiveAlternative(l Literal, e Expression) PrimitiveAlternative {
	return PrimitiveAlternative{l, e}
}

// Literal returns a literal pattern.
func (a PrimitiveAlternative) Literal() Literal {
	return a.literal
}

// Expression is an expression.
func (a PrimitiveAlternative) Expression() Expression {
	return a.expression
}

// VisitExpressions visits expressions.
func (a PrimitiveAlternative) VisitExpressions(f func(Expression) error) error {
	return a.expression.VisitExpressions(f)
}

// ConvertTypes converts types.
func (a PrimitiveAlternative) ConvertTypes(f func(types.Type) types.Type) PrimitiveAlternative {
	return PrimitiveAlternative{a.literal, a.expression.ConvertTypes(f)}
}

// RenameVariables renames variables.
func (a PrimitiveAlternative) RenameVariables(vs map[string]string) PrimitiveAlternative {
	return PrimitiveAlternative{a.literal, a.expression.RenameVariables(vs)}
}
