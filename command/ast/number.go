package ast

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

func (Number) isExpression() {}
func (Number) isLiteral()    {}
