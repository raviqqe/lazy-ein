package types

import "github.com/ein-lang/ein/command/debug"

func newTypeError(m string, i *debug.Information) error {
	return debug.NewError("TypeError", m, i)
}
