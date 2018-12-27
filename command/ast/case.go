package ast

// Case is a case expression.
type Case struct {
	expression         Expression
	alternatives       []Alternative
	defaultAlternative DefaultAlternative
}

// NewCase creates a case expression.
func NewCase(e Expression, as []Alternative, d DefaultAlternative) Case {
	return Case{e, as, d}
}

// NewCaseWithoutDefault creates a case expression.
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

// DefaultAlternative returns an expression.
func (c Case) DefaultAlternative() (DefaultAlternative, bool) {
	if c.defaultAlternative == (DefaultAlternative{}) {
		return DefaultAlternative{}, false
	}

	return c.defaultAlternative, true
}

// ConvertExpressions visits expressions.
func (c Case) ConvertExpressions(f func(Expression) Expression) node {
	as := make([]Alternative, 0, len(c.alternatives))

	for _, a := range c.alternatives {
		as = append(as, a.ConvertExpressions(f).(Alternative))
	}

	return f(
		Case{
			c.expression.ConvertExpressions(f).(Expression),
			as,
			c.defaultAlternative.ConvertExpressions(f).(DefaultAlternative),
		},
	)
}

func (Case) isExpression() {}
