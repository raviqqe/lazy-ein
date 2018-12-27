package ast

// Node is an AST node.
type Node interface {
	ConvertExpressions(func(Expression) Expression) Node
}
