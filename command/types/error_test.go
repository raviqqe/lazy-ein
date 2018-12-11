package types

import (
	"testing"

	"github.com/ein-lang/ein/command/debug"
	"github.com/stretchr/testify/assert"
)

func TestTypeErrorError(t *testing.T) {
	assert.Equal(
		t,
		"TypeError: you are wrong\nfoo.go:42:2049:\tfunc foo() {}",
		newTypeError(
			"you are wrong",
			debug.NewInformation("foo.go", 42, 2049, "func foo() {}"),
		).Error(),
	)
}
