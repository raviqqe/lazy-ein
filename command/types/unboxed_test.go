package types_test

import (
	"testing"

	"github.com/raviqqe/jsonxx/command/debug"
	"github.com/raviqqe/jsonxx/command/types"
	"github.com/stretchr/testify/assert"
)

func TestNewUnboxed(t *testing.T) {
	types.NewUnboxed(types.NewNumber(nil), nil)
}

func TestNewUnboxedPanic(t *testing.T) {
	assert.Panics(t, func() { types.NewUnboxed(types.NewUnboxed(types.NewNumber(nil), nil), nil) })
}

func TestUnboxedContent(t *testing.T) {
	assert.Equal(t, types.NewNumber(nil), types.NewUnboxed(types.NewNumber(nil), nil).Content())
}

func TestUnboxedUnify(t *testing.T) {
	assert.Nil(
		t,
		types.NewUnboxed(types.NewNumber(nil), nil).Unify(
			types.NewUnboxed(types.NewNumber(nil), nil),
		),
	)
}

func TestUnboxedUnifyError(t *testing.T) {
	assert.Error(t, types.NewUnboxed(types.NewNumber(nil), nil).Unify(types.NewNumber(nil)))
}

func TestUnboxedDebugInformation(t *testing.T) {
	i := debug.NewInformation("", 1, 1, "")
	assert.Equal(t, i, types.NewUnboxed(types.NewNumber(nil), i).DebugInformation())
}