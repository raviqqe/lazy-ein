package ast

import "github.com/raviqqe/stg/types"

// PrimitiveCase is a primitive case expression.
type PrimitiveCase struct {
	expression         Expression
	expressionType     types.Type
	alternatives       []PrimitiveAlternative
	defaultAlternative DefaultAlternative
}

// NewPrimitiveCase creates a primitive case expression.
func NewPrimitiveCase(e Expression, t types.Type, as []PrimitiveAlternative, a DefaultAlternative) PrimitiveCase {
	return PrimitiveCase{e, t, as, a}
}

// NewPrimitiveCaseWithoutDefault creates a primitive case expression.
func NewPrimitiveCaseWithoutDefault(e Expression, t types.Type, as []PrimitiveAlternative) PrimitiveCase {
	return PrimitiveCase{e, t, as, DefaultAlternative{}}
}

// Expression returns an expression.
func (c PrimitiveCase) Expression() Expression {
	return c.expression
}

// Type returns a type.
func (c PrimitiveCase) Type() types.Type {
	return c.expressionType
}

// Alternatives returns alternatives.
func (c PrimitiveCase) Alternatives() []PrimitiveAlternative {
	return c.alternatives
}

// DefaultAlternative returns a default alternative.
func (c PrimitiveCase) DefaultAlternative() (DefaultAlternative, bool) {
	if c.defaultAlternative == (DefaultAlternative{}) {
		return DefaultAlternative{}, false
	}

	return c.defaultAlternative, true
}

func (PrimitiveCase) isExpression() {}
