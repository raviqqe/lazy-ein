package ast

import "github.com/ein-lang/ein/command/core/types"

type abstractCase struct {
	argument           Expression
	argumentType       types.Type
	defaultAlternative DefaultAlternative
}

func newAbstractCase(e Expression, t types.Type, a DefaultAlternative) abstractCase {
	return abstractCase{e, t, a}
}

func (c abstractCase) Argument() Expression {
	return c.argument
}

func (c abstractCase) Type() types.Type {
	return c.argumentType
}

func (c abstractCase) DefaultAlternative() (DefaultAlternative, bool) {
	if c.defaultAlternative == (DefaultAlternative{}) {
		return DefaultAlternative{}, false
	}

	return c.defaultAlternative, true
}

func (c abstractCase) ConvertTypes(f func(types.Type) types.Type) abstractCase {
	d, ok := c.DefaultAlternative()

	if ok {
		d = c.defaultAlternative.ConvertTypes(f)
	}

	return abstractCase{c.argument.ConvertTypes(f), c.argumentType.ConvertTypes(f), d}
}

func (abstractCase) isExpression() {}
