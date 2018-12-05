package types_test

import (
	"testing"

	"github.com/raviqqe/jsonxx/command/types"
	"github.com/stretchr/testify/assert"
)

func TestFunctionUnify(t *testing.T) {
	assert.Nil(
		t,
		types.NewFunction(
			types.NewNumber(debugInformation),
			types.NewNumber(debugInformation),
			debugInformation,
		).Unify(
			types.NewFunction(
				types.NewNumber(debugInformation),
				types.NewNumber(debugInformation),
				debugInformation,
			),
		),
	)
}

func TestFunctionUnifyError(t *testing.T) {
	for _, ts := range [][2]types.Type{
		{
			types.NewFunction(
				types.NewNumber(debugInformation),
				types.NewNumber(debugInformation),
				debugInformation,
			),
			types.NewNumber(debugInformation),
		},
		{
			types.NewFunction(
				types.NewFunction(
					types.NewNumber(debugInformation),
					types.NewNumber(debugInformation),
					debugInformation,
				),
				types.NewNumber(debugInformation),
				debugInformation,
			),
			types.NewFunction(
				types.NewNumber(debugInformation),
				types.NewNumber(debugInformation),
				debugInformation,
			),
		},
	} {
		assert.Error(t, ts[0].Unify(ts[1]))
	}
}

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
