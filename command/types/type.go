package types

import (
	coretypes "github.com/raviqqe/lazy-ein/command/core/types"
	"github.com/raviqqe/lazy-ein/command/debug"
)

// Type is a type.
type Type interface {
	TypeVisitor
	Unify(Type) ([]Equation, error)
	SubstituteVariable(Variable, Type) Type
	DebugInformation() *debug.Information
	ToCore() coretypes.Type
}
