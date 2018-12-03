package types

import "github.com/raviqqe/jsonxx/command/debug"

// Function is a function type.
type Function struct {
	argument         Type
	result           Type
	debugInformation *debug.Information
}

// NewFunction creates a new function type.
func NewFunction(a, r Type, i *debug.Information) Function {
	return Function{a, r, i}
}

// Unify unifies itself with another type.
func (Function) Unify(t Type) error {
	if _, ok := t.(Function); ok {
		return nil
	}

	return newTypeError("not a function", t.DebugInformation())
}

// DebugInformation returns debug information.
func (f Function) DebugInformation() *debug.Information {
	return f.debugInformation
}
