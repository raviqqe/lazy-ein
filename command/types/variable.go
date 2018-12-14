package types

import (
	coretypes "github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/debug"
)

// Variable is a type variable.
type Variable struct {
	inferredType     Type
	debugInformation *debug.Information
}

// NewVariable creates a new variable.
func NewVariable(i *debug.Information) *Variable {
	return &Variable{nil, i}
}

// Unify unifies itself with another type.
func (v *Variable) Unify(t Type) error {
	if v.inferredType != nil {
		return v.inferredType.Unify(t)
	}

	v.inferredType = t

	return nil
}

// DebugInformation returns debug information.
func (v Variable) DebugInformation() *debug.Information {
	return v.debugInformation
}

// ToCore returns a type in the core language.
func (Variable) ToCore() coretypes.Type {
	panic("unreachable")
}
