package ast

import "github.com/ein-lang/ein/command/core/types"

type abstractCase struct {
	argument           Expression
	defaultAlternative DefaultAlternative
}

func newAbstractCase(e Expression, a DefaultAlternative) abstractCase {
	return abstractCase{e, a}
}

func (c abstractCase) Argument() Expression {
	return c.argument
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

	return abstractCase{c.argument.ConvertTypes(f), d}
}

func (c abstractCase) RenameVariables(vs map[string]string) abstractCase {
	d, ok := c.DefaultAlternative()

	if ok {
		d = c.defaultAlternative.RenameVariables(vs)
	}

	return abstractCase{c.argument.RenameVariables(vs), d}
}

func (abstractCase) isExpression() {}
