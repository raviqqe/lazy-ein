package debug_test

import (
	"testing"

	"github.com/ein-lang/ein/command/debug"
	"github.com/stretchr/testify/assert"
)

func TestErrorDebugInformation(t *testing.T) {
	i := debug.NewInformation("foo.go", 42, 2049, "func foo() {}")
	assert.Equal(t, i, debug.NewError("MyError", "you are wrong", i).DebugInformation())
}

func TestErrorError(t *testing.T) {
	assert.Equal(
		t,
		"MyError: you are wrong\nfoo.go:42:2049:\tfunc foo() {}",
		debug.NewError(
			"MyError",
			"you are wrong",
			debug.NewInformation("foo.go", 42, 2049, "func foo() {}"),
		).Error(),
	)
}
