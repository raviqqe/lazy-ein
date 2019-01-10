package ast

import "github.com/ein-lang/ein/command/core/types"

// PrimitiveCase is a primitive case expression.
type PrimitiveCase struct {
	abstractCase
	alternatives []PrimitiveAlternative
}

// NewPrimitiveCase creates a primitive case expression.
func NewPrimitiveCase(e Expression, t types.Type, as []PrimitiveAlternative, a DefaultAlternative) PrimitiveCase {
	return PrimitiveCase{newAbstractCase(e, t, a), as}
}

// NewPrimitiveCaseWithoutDefault creates a primitive case expression.
func NewPrimitiveCaseWithoutDefault(e Expression, t types.Type, as []PrimitiveAlternative) PrimitiveCase {
	return PrimitiveCase{newAbstractCase(e, t, DefaultAlternative{}), as}
}

// Alternatives returns alternatives.
func (c PrimitiveCase) Alternatives() []PrimitiveAlternative {
	return c.alternatives
}

// ConvertTypes converts types.
func (c PrimitiveCase) ConvertTypes(f func(types.Type) types.Type) Expression {
	as := make([]PrimitiveAlternative, 0, len(c.alternatives))

	for _, a := range c.alternatives {
		as = append(as, a.ConvertTypes(f))
	}

	return PrimitiveCase{c.abstractCase.ConvertTypes(f), as}
}
