package types

import (
	"fmt"

	"github.com/ein-lang/ein/command/debug"
)

// TypeError is a type error.
type TypeError struct {
	message          string
	debugInformation *debug.Information
}

func newTypeError(m string, i *debug.Information) error {
	return TypeError{m, i}
}

func (e TypeError) Error() string {
	return fmt.Sprintf("TypeError: %v\n%v", e.message, e.debugInformation)
}
