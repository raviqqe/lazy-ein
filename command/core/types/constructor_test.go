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
