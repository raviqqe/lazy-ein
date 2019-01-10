package ast

import "github.com/ein-lang/ein/command/core/types"

// Argument is a argument.
type Argument struct {
	name string
	typ  types.Type
}

// NewArgument creates an argument.
func NewArgument(n string, t types.Type) Argument {
	return Argument{n, t}
}

// Name returns a name.
func (a Argument) Name() string {
	return a.name
}

// Type returns a type.
func (a Argument) Type() types.Type {
	return a.typ
}

// ConvertTypes converts types.
func (a Argument) ConvertTypes(f func(types.Type) types.Type) Argument {
	return Argument{a.name, a.typ.ConvertTypes(f)}
}
