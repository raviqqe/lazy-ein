package types_test

import (
	"testing"

	"github.com/raviqqe/jsonxx/command/types"
	"github.com/stretchr/testify/assert"
)

func TestVariableUnify(t *testing.T) {
	for _, ts := range [][2]types.Type{
		{types.NewVariable(debugInformation), types.NewVariable(debugInformation)},
		{types.NewVariable(debugInformation), types.NewNumber(debugInformation)},
	} {
		assert.Nil(t, ts[0].Unify(ts[1]))
	}
}

func TestVariableUnifyWithUnifiedVariables(t *testing.T) {
	tt := types.NewVariable(debugInformation)
	assert.Nil(t, tt.Unify(types.NewNumber(debugInformation)))
	assert.Nil(t, tt.Unify(types.NewNumber(debugInformation)))
}

func TestVariableDebugInformation(t *testing.T) {
	assert.Equal(t, debugInformation, types.NewVariable(debugInformation).DebugInformation())
}
