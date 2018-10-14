package ast

import "llvm.org/llvm/bindings/go/llvm"

// Argument is a argument.
type Argument struct {
	name string
	typ  llvm.Type
}

// NewArgument creates a new argument.
func NewArgument(n string, t llvm.Type) Argument {
	return Argument{n, t}
}

// Name returns a name.
func (v Argument) Name() string {
	return v.name
}

// Type returns a type.
func (v Argument) Type() llvm.Type {
	return v.typ
}
