package types

import (
	"testing"

	"github.com/ein-lang/ein/command/debug"
)

func TestNewTypeError(t *testing.T) {
	newTypeError("you are wrong", debug.NewInformation("foo.go", 42, 2049, "func foo() {}"))
}
