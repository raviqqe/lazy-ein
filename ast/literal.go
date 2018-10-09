package ast

// Literal is a literal type.
type Literal interface {
	Expression
	isLiteral()
}
