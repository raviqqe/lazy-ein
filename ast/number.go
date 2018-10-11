package ast

// Number is a number literal.
type Number float64

// NewNumber creates a new number.
func NewNumber(n float64) Number {
	return Number(n)
}

// Value returns a value of a number.
func (n Number) Value() float64 {
	return float64(n)
}

func (Number) isExpression() {}
func (Number) isLiteral()    {}
