package parse

import (
	"testing"

	"github.com/raviqqe/jsonxx/command/debug"
	"github.com/stretchr/testify/assert"
)

func TestErrorError(t *testing.T) {
	assert.Equal(t, "foo", newError("foo", debug.NewInformation("", 0, 0, "")).Error())
}

func TestErrorDebugInformation(t *testing.T) {
	i := debug.NewInformation("", 0, 0, "")
	assert.Equal(t, i, newError("foo", i).DebugInformation())
}
