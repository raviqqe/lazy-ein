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
	assert.Panics(t, func() { NewBoxed(NewBoxed(NewFunction([]Type{NewFloat64()}, NewFloat64()))) })
}

func TestBoxedString(t *testing.T) {
	assert.Equal(t, "Boxed(Float64)", NewBoxed(NewFloat64()).String())
}

func TestBoxedEqual(t *testing.T) {
	assert.True(t, NewBoxed(NewFloat64()).equal(NewBoxed(NewFloat64())))
	assert.False(t, NewBoxed(NewFloat64()).equal(NewFloat64()))
	assert.False(
		t,
		NewBoxed(NewFloat64()).equal(
			NewBoxed(NewAlgebraic([]Constructor{NewConstructor("Nil", nil)})),
		),
	)
}
