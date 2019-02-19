package ast

import "github.com/ein-lang/ein/command/core/types"

// PrimitiveCase is a primitive case expression.
type PrimitiveCase struct {
	abstractCase
	typ          types.Primitive
	alternatives []PrimitiveAlternative
}

// NewPrimitiveCase creates a primitive case expression.
func NewPrimitiveCase(e Expression, t types.Primitive, as []PrimitiveAlternative, a DefaultAlternative) PrimitiveCase {
	return PrimitiveCase{newAbstractCase(e, a), t, as}
}

// NewPrimitiveCaseWithoutDefault creates a primitive case expression.
func NewPrimitiveCaseWithoutDefault(e Expression, t types.Primitive, as []PrimitiveAlternative) PrimitiveCase {
	return PrimitiveCase{newAbstractCase(e, DefaultAlternative{}), t, as}
}

// Type is a type.
func (c PrimitiveCase) Type() types.Primitive {
	return c.typ
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

	return PrimitiveCase{
		c.abstractCase.ConvertTypes(f),
		c.typ.ConvertTypes(f).(types.Primitive),
		as,
	}
}

// RenameVariables renames variables.
func (c PrimitiveCase) RenameVariables(vs map[string]string) Expression {
	as := make([]PrimitiveAlternative, 0, len(c.alternatives))

	for _, a := range c.alternatives {
		as = append(as, a.RenameVariables(vs))
	}

	return PrimitiveCase{c.abstractCase.RenameVariables(vs), c.typ, as}
}
