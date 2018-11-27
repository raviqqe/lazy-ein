package ast

// DefaultAlternative is a default alternative.
type DefaultAlternative struct {
	variable   string
	expression Expression
}

// NewDefaultAlternative creates a default alternative.
func NewDefaultAlternative(s string, e Expression) DefaultAlternative {
	return DefaultAlternative{s, e}
}

// Variable returns a bound variable.
func (a DefaultAlternative) Variable() string {
	return a.variable
}

// Expression is an expression.
func (a DefaultAlternative) Expression() Expression {
	return a.expression
}
