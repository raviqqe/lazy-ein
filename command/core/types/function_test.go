package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionString(t *testing.T) {
	assert.Equal(
		t,
		"Function([Float64,Float64],Float64)",
		NewFunction([]Type{NewFloat64(), NewFloat64()}, NewFloat64()).String(),
	)
}
