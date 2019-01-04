package ast

import "github.com/ein-lang/ein/command/core/types"

// TypeDefinition is a constructor definition.
type TypeDefinition struct {
	name string
	typ  types.Type
}

// NewTypeDefinition creates a constructor definition.
func NewTypeDefinition(n string, t types.Type) TypeDefinition {
	return TypeDefinition{n, t}
}

// Name returns a name.
func (d TypeDefinition) Name() string {
	return d.name
}

// Type returns a type.
func (d TypeDefinition) Type() types.Type {
	return d.typ
}
