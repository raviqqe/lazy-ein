package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstructorString(t *testing.T) {
	assert.Equal(t, "c", NewConstructor().String())
	assert.Equal(
		t,
		"c(f64,f64)",
		NewConstructor(NewFloat64(), NewFloat64()).String(),
	)
}

func TestConstructorEqual(t *testing.T) {
	for _, c := range []Constructor{
		NewConstructor(),
		NewConstructor(NewFloat64()),
	} {
		assert.True(t, c.equal(c))
	}

	for _, cs := range [][2]Constructor{
		{
			NewConstructor(),
			NewConstructor(NewFloat64()),
		},
		{
			NewConstructor(NewFloat64()),
			NewConstructor(NewFunction([]Type{NewFloat64()}, NewFloat64())),
		},
	} {
		assert.False(t, cs[0].equal(cs[1]))
	}
}
