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
