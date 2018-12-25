package types_test

import (
	"testing"

	"github.com/ein-lang/ein/command/debug"
	"github.com/ein-lang/ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestNewUnboxed(t *testing.T) {
	types.NewUnboxed(types.NewNumber(nil), nil)
}

func TestNewUnboxedPanic(t *testing.T) {
	assert.Panics(t, func() { types.NewUnboxed(types.NewUnboxed(types.NewNumber(nil), nil), nil) })
	assert.Panics(t, func() { types.NewUnboxed(types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil), nil) })
}

func TestUnboxedContent(t *testing.T) {
	assert.Equal(t, types.NewNumber(nil), types.NewUnboxed(types.NewNumber(nil), nil).Content())
}

func TestUnboxedUnify(t *testing.T) {
	es, err := types.NewUnboxed(types.NewNumber(nil), nil).Unify(
		types.NewUnboxed(types.NewNumber(nil), nil),
	)

	assert.Equal(t, []types.Equation(nil), es)
	assert.Nil(t, err)
}

func TestUnboxedUnifyError(t *testing.T) {
	_, err := types.NewUnboxed(types.NewNumber(nil), nil).Unify(types.NewNumber(nil))
	assert.Error(t, err)
}

func TestUnboxedDebugInformation(t *testing.T) {
	i := debug.NewInformation("", 1, 1, "")
	assert.Equal(t, i, types.NewUnboxed(types.NewNumber(nil), i).DebugInformation())
}
