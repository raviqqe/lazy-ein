package debug

import "fmt"

// Error is an error with debug information.
type Error struct {
	kind, message    string
	debugInformation *Information
}

// NewError creates an error.
func NewError(k, m string, i *Information) Error {
	return Error{k, m, i}
}

// DebugInformation returns debug information.
func (e Error) DebugInformation() *Information {
	return e.debugInformation
}

func (e Error) Error() string {
	return fmt.Sprintf("%v: %v", e.kind, e.message)
}
