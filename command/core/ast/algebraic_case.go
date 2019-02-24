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

// VisitExpressions visits expressions.
func (c AlgebraicCase) VisitExpressions(f func(Expression) error) error {
	for _, a := range c.alternatives {
		if err := a.VisitExpressions(f); err != nil {
			return err
		}
	}

	if d, ok := c.DefaultAlternative(); ok {
		if err := d.VisitExpressions(f); err != nil {
			return err
		}
	}

	return f(c)
}

// ConvertTypes converts types.
func (c AlgebraicCase) ConvertTypes(f func(types.Type) types.Type) Expression {
	as := make([]AlgebraicAlternative, 0, len(c.alternatives))

	for _, a := range c.alternatives {
		as = append(as, a.ConvertTypes(f))
	}

	return AlgebraicCase{c.abstractCase.ConvertTypes(f), as}
}

// RenameVariables renames variables.
func (c AlgebraicCase) RenameVariables(vs map[string]string) Expression {
	as := make([]AlgebraicAlternative, 0, len(c.alternatives))

	for _, a := range c.alternatives {
		as = append(as, a.RenameVariables(vs))
	}

	return AlgebraicCase{c.abstractCase.RenameVariables(vs), as}
}
