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
func (v Argument) Name() string {
	return v.name
}

// Type returns a type.
func (v Argument) Type() types.Type {
	return v.typ
}
