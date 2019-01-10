package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstructorString(t *testing.T) {
	assert.Equal(t, "Constructor(Nil)", NewConstructor("Nil", nil).String())
	assert.Equal(
		t,
		"Constructor(Cons,[Float64,Float64])",
		NewConstructor("Cons", []Type{NewFloat64(), NewFloat64()}).String(),
	)
}

func TestConstructorEqual(t *testing.T) {
	for _, c := range []Constructor{
		NewConstructor("Null", nil),
		NewConstructor("Cons", []Type{NewFloat64()}),
	} {
		assert.True(t, c.equal(c))
	}

	for _, cs := range [][2]Constructor{
		{
			NewConstructor("Nil", nil),
			NewConstructor("Null", nil),
		},
		{
			NewConstructor("Cons", []Type{NewFloat64()}),
			NewConstructor("Cons", []Type{NewBoxed(NewFloat64())}),
		},
		{
			NewConstructor("Cons", []Type{NewFloat64()}),
			NewConstructor("Cons", []Type{NewFloat64(), NewFloat64()}),
		},
	} {
		assert.False(t, cs[0].equal(cs[1]))
	}
}
