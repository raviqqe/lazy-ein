package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {
	assert.True(t, Equal(NewFloat64(), NewFloat64()))
	assert.False(t, Equal(NewFloat64(), NewBoxed(NewFloat64())))
}
