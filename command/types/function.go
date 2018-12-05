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
func (f Function) Unify(t Type) error {
	ff, ok := t.(Function)

	if !ok {
		return newTypeError("not a function", t.DebugInformation())
	} else if err := ff.argument.Unify(f.argument); err != nil {
		return err
	}

	return f.result.Unify(ff.result)
}

// DebugInformation returns debug information.
func (f Function) DebugInformation() *debug.Information {
	return f.debugInformation
}
