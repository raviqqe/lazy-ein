package ast

import "github.com/ein-lang/ein/command/types"

// Alternative is an alternative.
type Alternative struct {
	pattern, expression Expression
}

// NewAlternative creates an alternative.
func NewAlternative(p, e Expression) Alternative {
	return Alternative{p, e}
}

// Pattern returns a pattern.
func (a Alternative) Pattern() Expression {
	return a.pattern
}

// Expression is an expression.
func (a Alternative) Expression() Expression {
	return a.expression
}

// ConvertExpressions converts expressions.
func (a Alternative) ConvertExpressions(f func(Expression) Expression) Alternative {
	return NewAlternative(a.pattern, a.expression.ConvertExpressions(f).(Expression))
}

// VisitTypes visits types.
func (a Alternative) VisitTypes(f func(types.Type) error) error {
	if err := a.pattern.VisitTypes(f); err != nil {
		return err
	}

	return a.expression.VisitTypes(f)
}
