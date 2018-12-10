package types_test

import (
	"testing"

	"github.com/raviqqe/jsonxx/command/debug"
	"github.com/raviqqe/jsonxx/command/types"
	"github.com/stretchr/testify/assert"
)

func TestFunctionUnify(t *testing.T) {
	assert.Nil(
		t,
		types.NewFunction(
			types.NewNumber(nil),
			types.NewNumber(nil),
			nil,
		).Unify(types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil)),
	)
}

func TestFunctionUnifyError(t *testing.T) {
	for _, ts := range [][2]types.Type{
		{
			types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
			types.NewNumber(nil),
		},
		{
			types.NewFunction(
				types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
				types.NewNumber(nil),
				nil,
			),
			types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
		},
	} {
		assert.Error(t, ts[0].Unify(ts[1]))
	}
}

func TestFunctionDebugInformation(t *testing.T) {
	assert.Equal(
		t,
		(*debug.Information)(nil),
		types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil).DebugInformation(),
	)
}
