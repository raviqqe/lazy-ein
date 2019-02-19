package ast

import "github.com/ein-lang/ein/command/types"

// DefaultAlternative is a default alternative.
type DefaultAlternative struct {
	variable   string
	expression Expression
}

// NewDefaultAlternative creates a default alternative.
func NewDefaultAlternative(s string, e Expression) DefaultAlternative {
	return DefaultAlternative{s, e}
}

// Variable returns a bound variable.
func (a DefaultAlternative) Variable() string {
	return a.variable
}

// Expression is an expression.
func (a DefaultAlternative) Expression() Expression {
	return a.expression
}

// ConvertExpressions converts expressions.
func (a DefaultAlternative) ConvertExpressions(f func(Expression) Expression) Node {
	return NewDefaultAlternative(a.variable, a.expression.ConvertExpressions(f).(Expression))
}

// VisitTypes visits types.
func (a DefaultAlternative) VisitTypes(f func(types.Type) error) error {
	return a.expression.VisitTypes(f)
}
