package ast

import "github.com/ein-lang/ein/command/core/types"

// ConstructorApplication is a constructor application.
type ConstructorApplication struct {
	name      string
	arguments []Atom
}

// NewConstructorApplication creates a constructor application.
func NewConstructorApplication(n string, as []Atom) ConstructorApplication {
	return ConstructorApplication{n, as}
}

// Name returns a name.
func (a ConstructorApplication) Name() string {
	return a.name
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
