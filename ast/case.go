package ast

import "github.com/raviqqe/stg/types"

// Case is a case expression.
type Case struct {
	expression         Expression
	expressionType     types.Type
	alternatives       Alternatives
	defaultAlternative DefaultAlternative
}

// NewCase creates a case expression.
func NewCase(e Expression, t types.Type, as Alternatives, a DefaultAlternative) Case {
	return Case{e, t, as, a}
}

// NewCaseWithoutDefault creates a case expression.
func NewCaseWithoutDefault(e Expression, t types.Type, as Alternatives) Case {
	return Case{e, t, as, DefaultAlternative{}}
}

// Expression returns an expression.
func (c Case) Expression() Expression {
	return c.expression
}

// ExpressionType returns an expression type.
func (c Case) ExpressionType() types.Type {
	return c.expressionType
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

func (c Case) isExpression() {}
