package ast

import "github.com/ein-lang/ein/command/core/types"

// Expression is an expression.
type Expression interface {
	isExpression()
	VisitExpressions(func(Expression) error) error
	ConvertTypes(func(types.Type) types.Type) Expression
	RenameVariables(map[string]string) Expression
}
