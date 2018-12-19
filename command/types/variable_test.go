package types_test

import (
	"testing"

	"github.com/ein-lang/ein/command/debug"
	"github.com/ein-lang/ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestVariableUnify(t *testing.T) {
	v, n := types.NewVariable(0, nil), types.NewNumber(nil)
	es, err := v.Unify(n)

	assert.Equal(t, []types.Equation{types.NewEquation(v, n)}, es)
	assert.Nil(t, err)
}

func TestVariableDebugInformation(t *testing.T) {
	assert.Equal(t, (*debug.Information)(nil), types.NewVariable(0, nil).DebugInformation())
}
