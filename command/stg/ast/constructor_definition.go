package ast

import "github.com/raviqqe/jsonxx/command/stg/types"

// ConstructorDefinition is a constructor definition.
type ConstructorDefinition struct {
	name  string
	typ   types.Algebraic
	index int
}

// NewConstructorDefinition creates a constructor definition.
func NewConstructorDefinition(n string, t types.Algebraic, i int) ConstructorDefinition {
	return ConstructorDefinition{n, t, i}
}

// Name returns a name.
func (d ConstructorDefinition) Name() string {
	return d.name
}

// Type returns a type.
func (d ConstructorDefinition) Type() types.Algebraic {
	return d.typ
}

// Index returns an index.
func (d ConstructorDefinition) Index() int {
	return d.index
}
