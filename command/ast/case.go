package ast

import "github.com/ein-lang/ein/command/types"

// Case is a case expression.
type Case struct {
	argument           Expression
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

// Argument returns an argument.
func (c Case) Argument() Expression {
	return c.argument
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

// ConvertExpressions converts expressions.
func (c Case) ConvertExpressions(f func(Expression) Expression) Node {
	as := make([]Alternative, 0, len(c.alternatives))

	for _, a := range c.alternatives {
		as = append(as, a.ConvertExpressions(f).(Alternative))
	}

	d, ok := c.DefaultAlternative()

	if ok {
		d = c.defaultAlternative.ConvertExpressions(f).(DefaultAlternative)
	}

	return f(Case{c.argument.ConvertExpressions(f).(Expression), c.typ, as, d})
}

// VisitTypes visits types.
func (c Case) VisitTypes(f func(types.Type) error) error {
	if err := c.argument.VisitTypes(f); err != nil {
		return err
	} else if err := c.typ.VisitTypes(f); err != nil {
		return err
	}

	for _, a := range c.alternatives {
		if err := a.VisitTypes(f); err != nil {
			return err
		}
	}

	if d, ok := c.DefaultAlternative(); ok {
		return d.VisitTypes(f)
	}

	return nil
}

func (Case) isExpression() {}
