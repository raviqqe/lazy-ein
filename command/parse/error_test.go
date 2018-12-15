package parse

import (
	"testing"

	"github.com/ein-lang/ein/command/debug"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	assert.Error(t, newError("foo", debug.NewInformation("", 0, 0, "")))
}
