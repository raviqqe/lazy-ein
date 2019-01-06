package ast

import "github.com/ein-lang/ein/command/types"

// Node is an AST node.
type Node interface {
	types.TypeVisitor
	ConvertExpressions(func(Expression) Expression) Node
}
