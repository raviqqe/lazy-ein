package ast

import "github.com/raviqqe/jsonxx/command/core/types"

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
