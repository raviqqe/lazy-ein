package ast

// PrimitiveAlternative is a primitive alternative.
type PrimitiveAlternative struct {
	literal    Literal
	expression Expression
}

// NewPrimitiveAlternative creates a new primitive alternative.
func NewPrimitiveAlternative(l Literal, e Expression) PrimitiveAlternative {
	return PrimitiveAlternative{l, e}
}

// Literal returns a literal pattern.
func (a PrimitiveAlternative) Literal() Literal {
	return a.literal
}

// Expression is an expression.
func (a PrimitiveAlternative) Expression() Expression {
	return a.expression
}

func (a PrimitiveAlternative) isAlternative() {}
