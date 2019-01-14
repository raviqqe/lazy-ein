package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBoxed(t *testing.T) {
	NewBoxed(NewAlgebraic([]Constructor{NewConstructor([]Type{NewFloat64()})}))
}

func TestBoxedString(t *testing.T) {
	assert.Equal(
		t,
		"Boxed(Algebraic(Constructor(Float64)))",
		NewBoxed(NewAlgebraic([]Constructor{NewConstructor([]Type{NewFloat64()})})).String(),
	)
}

func TestBoxedEqual(t *testing.T) {
	a := NewAlgebraic([]Constructor{NewConstructor([]Type{NewFloat64()})})

	assert.True(t, NewBoxed(a).equal(NewBoxed(a)))
	assert.False(t, NewBoxed(a).equal(a))
	assert.False(
		t,
		NewBoxed(a).equal(
			NewBoxed(NewAlgebraic([]Constructor{NewConstructor(nil), NewConstructor(nil)})),
		),
	)
}
