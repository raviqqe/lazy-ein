package types

import "github.com/raviqqe/jsonxx/command/debug"

// Type is a type.
type Type interface {
	Unify(Type) error
	DebugInformation() *debug.Information
}
