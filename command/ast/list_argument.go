package ast

import "github.com/ein-lang/ein/command/types"

// ListArgument is a list argument.
type ListArgument struct {
	expression Expression
	expanded   bool
}

// NewListArgument creates an expanded expression.
func NewListArgument(e Expression, b bool) ListArgument {
	return ListArgument{e, b}
}

// Expression returns an expression.
func (a ListArgument) Expression() Expression {
	return a.expression
}

// Expanded returns true if the argument is expanded.
func (a ListArgument) Expanded() bool {
	return a.expanded
}

// ConvertExpressions visits expressions.
func (a ListArgument) ConvertExpressions(f func(Expression) Expression) Node {
	return ListArgument{a.expression.ConvertExpressions(f).(Expression), a.expanded}
}

// VisitTypes visits types.
func (a ListArgument) VisitTypes(f func(types.Type) error) error {
	return a.expression.VisitTypes(f)
}
