package ast

// Unboxed is an unboxed literal.
type Unboxed struct {
	content Literal
}

// NewUnboxed creates an unboxed literal.
func NewUnboxed(l Literal) Unboxed {
	return Unboxed{l}
}

// Content returns a content.
func (u Unboxed) Content() Literal {
	return u.content
}

// ConvertExpression visits expressions.
func (u Unboxed) ConvertExpression(func(Expression) Expression) node {
	return u
}

func (Unboxed) isExpression() {}
func (Unboxed) isLiteral()    {}
