package parse

import "github.com/ein-lang/ein/command/debug"

func newError(s string, i *debug.Information) error {
	return debug.NewError("ParseError", s, i)
}
