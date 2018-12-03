package types_test

import (
	"testing"

	"github.com/raviqqe/jsonxx/command/types"
	"github.com/stretchr/testify/assert"
)

func TestFunctionDebugInformation(t *testing.T) {
	assert.Equal(
		t,
		debugInformation,
		types.NewFunction(
			types.NewNumber(debugInformation),
			types.NewNumber(debugInformation),
			debugInformation,
		).DebugInformation(),
	)
}
