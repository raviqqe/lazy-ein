package stg

// Number is a number literal.
type Number float64

// NewNumber creates a new number.
func NewNumber(n float64) Number {
	return Number(n)
}

func (Number) isExpression() {}
func (Number) isLiteral()    {}
