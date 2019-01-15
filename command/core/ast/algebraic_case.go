package ast

import "github.com/ein-lang/ein/command/core/types"

// AlgebraicCase is an algebraic case expression.
type AlgebraicCase struct {
	abstractCase
	alternatives []AlgebraicAlternative
}

// NewAlgebraicCase creates an algebraic case expression.
func NewAlgebraicCase(e Expression, as []AlgebraicAlternative, a DefaultAlternative) AlgebraicCase {
	return AlgebraicCase{newAbstractCase(e, a), as}
}

// NewAlgebraicCaseWithoutDefault creates an algebraic case expression.
func NewAlgebraicCaseWithoutDefault(e Expression, as []AlgebraicAlternative) AlgebraicCase {
	return AlgebraicCase{newAbstractCase(e, DefaultAlternative{}), as}
}

// Alternatives returns alternatives.
func (c AlgebraicCase) Alternatives() []AlgebraicAlternative {
	return c.alternatives
}

// ConvertTypes converts types.
func (c AlgebraicCase) ConvertTypes(f func(types.Type) types.Type) Expression {
	as := make([]AlgebraicAlternative, 0, len(c.alternatives))

	for _, a := range c.alternatives {
		as = append(as, a.ConvertTypes(f))
	}

	return AlgebraicCase{c.abstractCase.ConvertTypes(f), as}
}
