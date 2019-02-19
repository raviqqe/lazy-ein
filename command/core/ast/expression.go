package ast

import "github.com/ein-lang/ein/command/core/types"

// Expression is an expression.
type Expression interface {
	isExpression()
	ConvertTypes(func(types.Type) types.Type) Expression
	RenameVariables(map[string]string) Expression
}
