package ast

// Float64 is a float64 literal.
type Float64 float64

// NewFloat64 creates a new float64 number.
func NewFloat64(n float64) Float64 {
	return Float64(n)
}

// Value returns a value.
func (n Float64) Value() float64 {
	return float64(n)
}

func (Float64) isExpression() {}
func (Float64) isLiteral()    {}
