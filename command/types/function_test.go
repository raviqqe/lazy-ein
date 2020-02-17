package types_test

import (
	"testing"

	"github.com/raviqqe/lazy-ein/command/debug"
	"github.com/raviqqe/lazy-ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestFunctionUnify(t *testing.T) {
	es, err := types.NewFunction(
		types.NewNumber(nil),
		types.NewNumber(nil),
		nil,
	).Unify(types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil))

	assert.Nil(t, err)
	assert.Equal(t, []types.Equation(nil), es)
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
		{
			types.NewFunction(
				types.NewNumber(nil),
				types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
				nil,
			),
			types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
		},
	} {
		_, err := ts[0].Unify(ts[1])
		assert.Error(t, err)
	}
}

func TestFunctionDebugInformation(t *testing.T) {
	assert.Equal(
		t,
		(*debug.Information)(nil),
		types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil).DebugInformation(),
	)
}
