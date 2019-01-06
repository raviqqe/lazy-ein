package types

import (
	coretypes "github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/debug"
)

// Unknown is an unknown type.
type Unknown struct {
	debugInformation *debug.Information
}

// NewUnknown creates an unknown type.
func NewUnknown(i *debug.Information) Unknown {
	return Unknown{i}
}

// Unify unifies itself with another type.
func (Unknown) Unify(Type) ([]Equation, error) {
	panic("unreachable")
}

// SubstituteVariable substitutes type variables.
func (Unknown) SubstituteVariable(Variable, Type) Type {
	panic("unreachable")
}

// DebugInformation returns debug information.
func (u Unknown) DebugInformation() *debug.Information {
	return u.debugInformation
}

// ToCore returns a type in the core language.
func (Unknown) ToCore() coretypes.Type {
	panic("unreachable")
}

// VisitTypes visits types.
func (u Unknown) VisitTypes(f func(Type) error) error {
	return f(u)
}

func (u Unknown) coreName() string {
	panic("unreachable")
}
