package ast

import "github.com/ein-lang/ein/command/core/types"

type abstractCase struct {
	expression         Expression
	expressionType     types.Type
	defaultAlternative DefaultAlternative
}

func newAbstractCase(e Expression, t types.Type, a DefaultAlternative) abstractCase {
	return abstractCase{e, t, a}
}

func (c abstractCase) Expression() Expression {
	return c.expression
}

func (c abstractCase) Type() types.Type {
	return c.expressionType
}

func (c abstractCase) DefaultAlternative() (DefaultAlternative, bool) {
	if c.defaultAlternative == (DefaultAlternative{}) {
		return DefaultAlternative{}, false
	}

	return c.defaultAlternative, true
}

func (c abstractCase) ConvertTypes(f func(types.Type) types.Type) abstractCase {
	return abstractCase{
		c.expression.ConvertTypes(f),
		c.expressionType.ConvertTypes(f),
		c.defaultAlternative.ConvertTypes(f),
	}
}

func (abstractCase) isExpression() {}
