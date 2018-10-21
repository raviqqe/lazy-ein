package ast

import "github.com/raviqqe/stg/types"

// Argument is a argument.
type Argument struct {
	name string
	typ  types.Type
}

// NewArgument creates a new argument.
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
