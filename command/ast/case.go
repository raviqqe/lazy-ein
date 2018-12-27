package ast

import "github.com/ein-lang/ein/command/types"

// Case is a case expression.
type Case struct {
	expression         Expression
	typ                types.Type
	alternatives       []Alternative
	defaultAlternative DefaultAlternative
}

// NewCase creates a case expression.
func NewCase(e Expression, t types.Type, as []Alternative, d DefaultAlternative) Case {
	return Case{e, t, as, d}
}

// NewCaseWithoutDefault creates a case expression.
func NewCaseWithoutDefault(e Expression, t types.Type, as []Alternative) Case {
	return Case{e, t, as, DefaultAlternative{}}
}

// Expression returns an expression.
func (c Case) Expression() Expression {
	return c.expression
}

// Type returns an expression.
func (c Case) Type() types.Type {
	return c.typ
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
func (c Case) ConvertExpressions(f func(Expression) Expression) Node {
	as := make([]Alternative, 0, len(c.alternatives))

	for _, a := range c.alternatives {
		as = append(as, a.ConvertExpressions(f).(Alternative))
	}

	d, ok := c.DefaultAlternative()

	if ok {
		d = c.defaultAlternative.ConvertExpressions(f).(DefaultAlternative)
	}

	return f(Case{c.expression.ConvertExpressions(f).(Expression), c.typ, as, d})
}

func (Case) isExpression() {}
