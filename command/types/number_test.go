package types_test

import (
	"testing"

	"github.com/raviqqe/jsonxx/command/types"
	"github.com/stretchr/testify/assert"
)

func TestNumberUnify(t *testing.T) {
	assert.Nil(t, types.NewNumber(debugInformation).Unify(types.NewNumber(debugInformation)))
}

func TestNumberUnifyError(t *testing.T) {
	assert.Error(
		t,
		types.NewNumber(debugInformation).Unify(
			types.NewFunction(
				types.NewNumber(debugInformation),
				types.NewNumber(debugInformation),
				debugInformation,
			),
		),
	)
}

func TestNumberDebugInformation(t *testing.T) {
	assert.Equal(t, debugInformation, types.NewNumber(debugInformation).DebugInformation())
}
