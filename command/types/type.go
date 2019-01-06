package types

import (
	coretypes "github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/debug"
)

// Type is a type.
type Type interface {
	Unify(Type) ([]Equation, error)
	SubstituteVariable(Variable, Type) Type
	DebugInformation() *debug.Information
	ToCore() (coretypes.Type, error)
	coreName() (string, error)
}
