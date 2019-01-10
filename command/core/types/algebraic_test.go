package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlgebraicString(t *testing.T) {
	assert.Equal(
		t,
		"Algebraic([Constructor(Nil),Constructor(Cons)])",
		NewAlgebraic(
			[]Constructor{NewConstructor("Nil", nil), NewConstructor("Cons", nil)},
		).String(),
	)
}

func TestAlgebraicEqual(t *testing.T) {
	a := NewAlgebraic([]Constructor{NewConstructor("Foo", nil)})
	assert.True(t, a.equal(a))

	for _, as := range [][2]Type{
		{a, NewFloat64()},
		{a, NewAlgebraic([]Constructor{NewConstructor("Bar", nil)})},
		{a, NewAlgebraic([]Constructor{NewConstructor("Foo", nil), NewConstructor("Bar", nil)})},
	} {
		assert.False(t, as[0].equal(as[1]))
	}
}
