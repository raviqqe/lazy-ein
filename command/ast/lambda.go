package ast

import "github.com/ein-lang/ein/command/types"

// Lambda is a lambda.
type Lambda struct {
	arguments  []string
	expression Expression
}

// NewLambda creates a lambda.
func NewLambda(as []string, e Expression) Lambda {
	return Lambda{as, e}
}

// Arguments returns arguments.
func (l Lambda) Arguments() []string {
	return l.arguments
}

// Expression returns an expression.
func (l Lambda) Expression() Expression {
	return l.expression
}

// ConvertExpressions visits expressions.
func (l Lambda) ConvertExpressions(f func(Expression) Expression) Node {
	return f(NewLambda(l.arguments, l.expression.ConvertExpressions(f).(Expression)))
}

// VisitTypes visits types.
func (l Lambda) VisitTypes(f func(types.Type) error) error {
	return l.expression.VisitTypes(f)
}

func (Lambda) isExpression() {}
