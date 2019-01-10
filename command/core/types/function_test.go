package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFunctionPanic(t *testing.T) {
	assert.Panics(
		t,
		func() {
			NewFunction([]Type{NewFloat64()}, NewFunction([]Type{NewFloat64()}, NewFloat64()))
		},
	)
}

func TestFunctionString(t *testing.T) {
	assert.Equal(
		t,
		"Function([Float64,Float64],Float64)",
		NewFunction([]Type{NewFloat64(), NewFloat64()}, NewFloat64()).String(),
	)
}

func TestFunctionEqual(t *testing.T) {
	a := NewFunction([]Type{NewFloat64()}, NewFloat64())
	assert.True(t, a.equal(a))

	for _, as := range [][2]Type{
		{a, NewFloat64()},
		{a, NewFunction([]Type{NewBoxed(NewFloat64())}, NewFloat64())},
		{a, NewFunction([]Type{NewFloat64(), NewFloat64()}, NewFloat64())},
	} {
		assert.False(t, as[0].equal(as[1]))
	}
}
