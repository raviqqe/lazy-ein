package ast

// AlgebraicAlternative is an algebraic alternative.
type AlgebraicAlternative struct {
	constructorName string
	elementNames    []string
	expression      Expression
}

// NewAlgebraicAlternative creates an algebraic alternative.
func NewAlgebraicAlternative(c string, es []string, e Expression) AlgebraicAlternative {
	return AlgebraicAlternative{c, es, e}
}

// ConstructorName returns a constructor name.
func (a AlgebraicAlternative) ConstructorName() string {
	return a.constructorName
}

// ElementNames returns element names.
func (a AlgebraicAlternative) ElementNames() []string {
	return a.elementNames
}

// Expression is an expression.
func (a AlgebraicAlternative) Expression() Expression {
	return a.expression
}

func (AlgebraicAlternative) isAlternative() {}
