package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlgebraicString(t *testing.T) {
	assert.Equal(
		t,
		"Algebraic(Constructor,Constructor)",
		NewAlgebraic(NewConstructor(), NewConstructor()).String(),
	)
}

func TestAlgebraicEqual(t *testing.T) {
	a := NewAlgebraic(NewConstructor())
	assert.True(t, a.equal(a))

	for _, as := range [][2]Type{
		{a, NewFloat64()},
		{a, NewAlgebraic(NewConstructor(), NewConstructor())},
	} {
		assert.False(t, as[0].equal(as[1]))
	}
}
