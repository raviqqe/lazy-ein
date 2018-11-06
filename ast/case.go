package ast

// Case is a case expression.
type Case struct {
	expression         Expression
	alternatives       []Alternative
	defaultAlternative DefaultAlternative
}

// NewCase creates a new case expression.
func NewCase(e Expression, as []Alternative, a DefaultAlternative) Case {
	return Case{e, as, a}
}

// NewCaseWithoutDefault creates a new case expression.
func NewCaseWithoutDefault(e Expression, as []Alternative) Case {
	return Case{e, as, DefaultAlternative{}}
}

// Expression returns an expression.
func (c Case) Expression() Expression {
	return c.expression
}

// Alternatives returns alternatives.
func (c Case) Alternatives() []Alternative {
	return c.alternatives
}

// DefaultAlternative returns a default alternative.
func (c Case) DefaultAlternative() (DefaultAlternative, bool) {
	if c.defaultAlternative == (DefaultAlternative{}) {
		return DefaultAlternative{}, false
	}

	return c.defaultAlternative, true
}

func (c Case) isExpression() {}
