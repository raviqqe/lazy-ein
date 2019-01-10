package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat64Equal(t *testing.T) {
	assert.True(t, NewFloat64().equal(NewFloat64()))
	assert.False(t, NewFloat64().equal(NewBoxed(NewFloat64())))
}
