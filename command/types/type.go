package types

import "github.com/ein-lang/ein/command/debug"

// Type is a type.
type Type interface {
	Unify(Type) error
	DebugInformation() *debug.Information
}
