package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBoxed(t *testing.T) {
	NewBoxed(NewFloat64())
}

func TestNewBoxedPanic(t *testing.T) {
	assert.Panics(t, func() { NewBoxed(NewBoxed(NewFloat64())) })
}
