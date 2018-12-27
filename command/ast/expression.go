package ast

// Expression is an expression.
type Expression interface {
	Node
	isExpression()
}
