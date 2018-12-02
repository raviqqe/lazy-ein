package ast

// Literal is a literal
type Literal interface {
	Expression
	isLiteral()
}
