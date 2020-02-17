package types

import "github.com/raviqqe/lazy-ein/command/debug"

// NewTypeError creates a type error.
func NewTypeError(m string, i *debug.Information) error {
	return debug.NewError("TypeError", m, i)
}
