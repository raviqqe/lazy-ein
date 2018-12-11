package types_test

import (
	"testing"

	"github.com/ein-lang/ein/command/debug"
	"github.com/ein-lang/ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestNumberUnify(t *testing.T) {
	assert.Nil(t, types.NewNumber(nil).Unify(types.NewNumber(nil)))
}

func TestNumberUnifyError(t *testing.T) {
	assert.Error(
		t,
		types.NewNumber(nil).Unify(
			types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
		),
	)
}

func TestNumberDebugInformation(t *testing.T) {
	assert.Equal(t, (*debug.Information)(nil), types.NewNumber(nil).DebugInformation())
}
