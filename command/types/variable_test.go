package types_test

import (
	"testing"

	"github.com/raviqqe/jsonxx/command/debug"
	"github.com/raviqqe/jsonxx/command/types"
	"github.com/stretchr/testify/assert"
)

func TestVariableUnify(t *testing.T) {
	for _, ts := range [][2]types.Type{
		{types.NewVariable(nil), types.NewVariable(nil)},
		{types.NewVariable(nil), types.NewNumber(nil)},
	} {
		assert.Nil(t, ts[0].Unify(ts[1]))
	}
}

func TestVariableUnifyWithUnifiedVariables(t *testing.T) {
	tt := types.NewVariable(nil)
	assert.Nil(t, tt.Unify(types.NewNumber(nil)))
	assert.Nil(t, tt.Unify(types.NewNumber(nil)))
}

func TestVariableDebugInformation(t *testing.T) {
	assert.Equal(t, (*debug.Information)(nil), types.NewVariable(nil).DebugInformation())
}
