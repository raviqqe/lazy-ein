package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstructorString(t *testing.T) {
	assert.Equal(t, "Constructor", NewConstructor(nil).String())
	assert.Equal(
		t,
		"Constructor(Float64,Float64)",
		NewConstructor([]Type{NewFloat64(), NewFloat64()}).String(),
	)
}

func TestConstructorEqual(t *testing.T) {
	for _, c := range []Constructor{
		NewConstructor(nil),
		NewConstructor([]Type{NewFloat64()}),
	} {
		assert.True(t, c.equal(c))
	}

	for _, cs := range [][2]Constructor{
		{
			NewConstructor(nil),
			NewConstructor([]Type{NewFloat64()}),
		},
		{
			NewConstructor([]Type{NewFloat64()}),
			NewConstructor([]Type{NewBoxed(NewFloat64())}),
		},
	} {
		assert.False(t, cs[0].equal(cs[1]))
	}
}
