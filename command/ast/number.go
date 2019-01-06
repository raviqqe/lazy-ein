package ast

import "github.com/ein-lang/ein/command/types"

// Number is a number.
type Number struct {
	value float64
}

// NewNumber creates a number.
func NewNumber(n float64) Number {
	return Number{n}
}

// Value returns a value.
func (n Number) Value() float64 {
	return n.value
}

// ConvertExpressions visits expressions.
func (n Number) ConvertExpressions(f func(Expression) Expression) Node {
	return f(n)
}

// VisitTypes visits types.
func (n Number) VisitTypes(f func(types.Type) error) error {
	return nil
}

func (Number) isExpression() {}
func (Number) isLiteral()    {}
