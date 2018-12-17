package types

import (
	coretypes "github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/debug"
)

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
	} else if v, ok := t.(*Variable); ok {
		return v.Unify(n)
	}

	return NewTypeError("not a number", t.DebugInformation())
}

// DebugInformation returns debug information.
func (n Number) DebugInformation() *debug.Information {
	return n.debugInformation
}

// ToCore returns a type in the core language.
func (n Number) ToCore() (coretypes.Type, error) {
	return coretypes.NewBoxed(coretypes.NewFloat64()), nil
}
