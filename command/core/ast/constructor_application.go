package ast

import "github.com/ein-lang/ein/command/core/types"

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

// ConvertTypes converts types.
func (a ConstructorApplication) ConvertTypes(func(types.Type) types.Type) Expression {
	return a
}

func (ConstructorApplication) isExpression() {}
