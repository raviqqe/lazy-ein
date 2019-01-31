package ast

import (
	"testing"

	"github.com/ein-lang/ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestPatternsEqual(t *testing.T) {
	for _, e := range []Expression{
		NewNumber(42),
		NewList(types.NewUnknown(nil), nil),
		NewList(types.NewUnknown(nil), []Expression{NewNumber(42)}),
		NewList(types.NewUnknown(nil), []Expression{NewNumber(42), NewNumber(42)}),
		NewList(
			types.NewUnknown(nil),
			[]Expression{NewList(types.NewUnknown(nil), []Expression{NewNumber(42)})},
		),
	} {
		assert.True(t, PatternsEqual(e, e))
	}

	for _, es := range [][2]Expression{
		{NewNumber(42), NewNumber(0)},
		{NewList(types.NewUnknown(nil), nil), NewNumber(42)},
		{
			NewList(types.NewUnknown(nil), []Expression{NewNumber(42)}),
			NewList(types.NewUnknown(nil), []Expression{NewNumber(0)}),
		},
		{
			NewList(types.NewUnknown(nil), []Expression{NewNumber(42)}),
			NewList(types.NewUnknown(nil), []Expression{NewNumber(42), NewNumber(42)}),
		},
	} {
		assert.False(t, PatternsEqual(es[0], es[1]))
	}
}
