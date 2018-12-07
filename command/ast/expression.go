package ast

// Expression is an expression.
type Expression interface {
	node
	isExpression()
}
