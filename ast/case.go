package ast

// Case is a case expression.
type Case struct {
	expression         Expression
	alternatives       Alternatives
	defaultAlternative DefaultAlternative
}

// NewCase creates a case expression.
func NewCase(e Expression, as Alternatives, a DefaultAlternative) Case {
	return Case{e, as, a}
}

// NewCaseWithoutDefault creates a case expression.
func NewCaseWithoutDefault(e Expression, as Alternatives) Case {
	return Case{e, as, DefaultAlternative{}}
}

// Expression returns an expression.
func (c Case) Expression() Expression {
	return c.expression
}

// Alternatives returns alternatives.
func (c Case) Alternatives() []Alternative {
	return c.alternatives.toSlice()
}

// DefaultAlternative returns a default alternative.
func (c Case) DefaultAlternative() (DefaultAlternative, bool) {
	if c.defaultAlternative == (DefaultAlternative{}) {
		return DefaultAlternative{}, false
	}

	return c.defaultAlternative, true
}

// IsAlgebraic returns true if a case expression is algebraic.
func (c Case) IsAlgebraic() bool {
	_, ok := c.alternatives.(AlgebraicAlternatives)
	return ok
}

func (c Case) isExpression() {}
