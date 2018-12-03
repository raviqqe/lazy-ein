package types_test

import (
	"testing"

	"github.com/raviqqe/jsonxx/command/types"
	"github.com/stretchr/testify/assert"
)

func TestNumberUnify(t *testing.T) {
	assert.Nil(
		t,
		types.Type(types.NewNumber(debugInformation)).Unify(
			types.Type(types.NewNumber(debugInformation)),
		),
	)
}
