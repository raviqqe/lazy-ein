package types_test

import (
	"testing"

	"github.com/ein-lang/ein/command/debug"
	"github.com/ein-lang/ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestNumberUnify(t *testing.T) {
	_, err := types.NewNumber(nil).Unify(types.NewNumber(nil))
	assert.Nil(t, err)
}

func TestNumberUnifyError(t *testing.T) {
	_, err := types.NewNumber(nil).Unify(
		types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
	)
	assert.Error(t, err)
}

func TestNumberDebugInformation(t *testing.T) {
	assert.Equal(t, (*debug.Information)(nil), types.NewNumber(nil).DebugInformation())
}
