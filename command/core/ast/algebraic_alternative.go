package ast

import "github.com/ein-lang/ein/command/core/types"

// AlgebraicAlternative is an algebraic alternative.
type AlgebraicAlternative struct {
	constructor  Constructor
	elementNames []string
	expression   Expression
}

// NewAlgebraicAlternative creates an algebraic alternative.
func NewAlgebraicAlternative(c Constructor, es []string, e Expression) AlgebraicAlternative {
	return AlgebraicAlternative{c, es, e}
}

// Constructor returns a constructor.
func (a AlgebraicAlternative) Constructor() Constructor {
	return a.constructor
}

// ElementNames returns element names.
func (a AlgebraicAlternative) ElementNames() []string {
	return a.elementNames
}

// Expression is an expression.
func (a AlgebraicAlternative) Expression() Expression {
	return a.expression
}

// VisitExpressions visits expressions.
func (a AlgebraicAlternative) VisitExpressions(f func(Expression) error) error {
	return a.expression.VisitExpressions(f)
}

// ConvertTypes converts types.
func (a AlgebraicAlternative) ConvertTypes(f func(types.Type) types.Type) AlgebraicAlternative {
	return AlgebraicAlternative{a.constructor, a.elementNames, a.expression.ConvertTypes(f)}
}

// RenameVariables renames variables.
func (a AlgebraicAlternative) RenameVariables(vs map[string]string) AlgebraicAlternative {
	return AlgebraicAlternative{
		a.constructor,
		a.elementNames,
		a.expression.RenameVariables(removeVariables(vs, a.elementNames...)),
	}
}
