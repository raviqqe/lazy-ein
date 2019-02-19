package ast

import "github.com/ein-lang/ein/command/core/types"

// Declaration is a declaration.
type Declaration struct {
	name   string
	lambda LambdaDeclaration
}

// NewDeclaration creates a declaration.
func NewDeclaration(n string, l LambdaDeclaration) Declaration {
	return Declaration{n, l}
}

// Name returns a name.
func (d Declaration) Name() string {
	return d.name
}

// Lambda returns a lambda declaration.
func (d Declaration) Lambda() LambdaDeclaration {
	return d.lambda
}

// ConvertTypes converts types.
func (d Declaration) ConvertTypes(f func(types.Type) types.Type) Declaration {
	return Declaration{d.name, d.lambda.ConvertTypes(f)}
}
