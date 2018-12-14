package types

import (
	coretypes "github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/debug"
)

// Type is a type.
type Type interface {
	Unify(Type) error
	DebugInformation() *debug.Information
	ToCore() coretypes.Type
}
