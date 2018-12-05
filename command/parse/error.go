package parse

import "github.com/raviqqe/jsonxx/command/debug"

// Error is a parse error.
type Error struct {
	message          string
	debugInformation *debug.Information
}

func newError(s string, i *debug.Information) Error {
	return Error{s, i}
}

func (e Error) Error() string {
	return e.message
}

// DebugInformation returns debug information.
func (e Error) DebugInformation() *debug.Information {
	return e.debugInformation
}
