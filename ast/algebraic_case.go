package ast

import "github.com/raviqqe/stg/types"

// AlgebraicCase is an algebraic case expression.
type AlgebraicCase struct {
	expression         Expression
	expressionType     types.Type
	alternatives       []AlgebraicAlternative
	defaultAlternative DefaultAlternative
}

// NewAlgebraicCase creates an algebraic case expression.
func NewAlgebraicCase(e Expression, t types.Type, as []AlgebraicAlternative, a DefaultAlternative) AlgebraicCase {
	return AlgebraicCase{e, t, as, a}
}

// NewAlgebraicCaseWithoutDefault creates an algebraic case expression.
func NewAlgebraicCaseWithoutDefault(e Expression, t types.Type, as []AlgebraicAlternative) AlgebraicCase {
	return AlgebraicCase{e, t, as, DefaultAlternative{}}
}

// Expression returns an expression.
func (c AlgebraicCase) Expression() Expression {
	return c.expression
}

// Type returns a type.
func (c AlgebraicCase) Type() types.Type {
	return c.expressionType
}

// Alternatives returns alternatives.
func (c AlgebraicCase) Alternatives() []AlgebraicAlternative {
	return c.alternatives
}

// DefaultAlternative returns a default alternative.
func (c AlgebraicCase) DefaultAlternative() (DefaultAlternative, bool) {
	if c.defaultAlternative == (DefaultAlternative{}) {
		return DefaultAlternative{}, false
	}

	return c.defaultAlternative, true
}

func (AlgebraicCase) isExpression() {}
