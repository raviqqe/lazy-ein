package parse

import "github.com/raviqqe/lazy-ein/command/debug"

func newError(s string, i *debug.Information) error {
	return debug.NewError("SyntaxError", s, i)
}
