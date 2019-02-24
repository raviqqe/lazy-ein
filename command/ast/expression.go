package ast

import "github.com/ein-lang/ein/command/types"

// Expression is an expression.
type Expression interface {
	types.TypeVisitor
	ConvertExpressions(func(Expression) Expression) Expression
	isExpression()
}
