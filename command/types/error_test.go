package types

import (
	"testing"

	"github.com/raviqqe/lazy-ein/command/debug"
)

func TestNewTypeError(t *testing.T) {
	NewTypeError("you are wrong", debug.NewInformation("foo.go", 42, 2049, "func foo() {}"))
}
