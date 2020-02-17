package ast

import "github.com/raviqqe/lazy-ein/command/core/types"

// ConstructorApplication is a constructor application.
type ConstructorApplication struct {
	constructor Constructor
	arguments   []Atom
}

// NewConstructorApplication creates a constructor application.
func NewConstructorApplication(c Constructor, as []Atom) ConstructorApplication {
	return ConstructorApplication{c, as}
}

// Constructor returns a constructor.
func (a ConstructorApplication) Constructor() Constructor {
	return a.constructor
}

// Arguments returns arguments.
func (a ConstructorApplication) Arguments() []Atom {
	return a.arguments
}

// VisitExpressions visits expressions.
func (a ConstructorApplication) VisitExpressions(f func(Expression) error) error {
	return f(a)
}

// ConvertTypes converts types.
func (a ConstructorApplication) ConvertTypes(func(types.Type) types.Type) Expression {
	return a
}

// RenameVariables renames variables.
func (a ConstructorApplication) RenameVariables(vs map[string]string) Expression {
	as := make([]Atom, 0, len(a.arguments))

	for _, a := range a.arguments {
		as = append(as, a.RenameVariablesInAtom(vs))
	}

	return ConstructorApplication{a.constructor, as}
}

func (ConstructorApplication) isExpression() {}
