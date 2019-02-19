package ast

import "github.com/ein-lang/ein/command/types"

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

// ConvertExpressions converts expressions.
func (u Unboxed) ConvertExpressions(f func(Expression) Expression) Node {
	return f(u)
}

// VisitTypes visits types.
func (u Unboxed) VisitTypes(f func(types.Type) error) error {
	return u.content.VisitTypes(f)
}

func (Unboxed) isExpression() {}
func (Unboxed) isLiteral()    {}
