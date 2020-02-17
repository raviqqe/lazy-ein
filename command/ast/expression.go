package ast

import "github.com/raviqqe/lazy-ein/command/types"

// Expression is an expression.
type Expression interface {
	types.TypeVisitor
	ConvertExpressions(func(Expression) Expression) Expression
	isExpression()
}
