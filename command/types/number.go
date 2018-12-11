package types

import "github.com/ein-lang/ein/command/debug"

// Number is a number type.
type Number struct {
	debugInformation *debug.Information
}

// NewNumber creates a new number type.
func NewNumber(i *debug.Information) Number {
	return Number{i}
}

// Unify unifies itself with another type.
func (n Number) Unify(t Type) error {
	if _, ok := t.(Number); ok {
		return nil
	}

	return newTypeError("not a number", t.DebugInformation())
}

// DebugInformation returns debug information.
func (n Number) DebugInformation() *debug.Information {
	return n.debugInformation
}
