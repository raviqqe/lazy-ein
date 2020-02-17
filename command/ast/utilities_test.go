package ast

import (
	"testing"

	"github.com/raviqqe/lazy-ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestPatternsEqual(t *testing.T) {
	for _, e := range []Expression{
		NewVariable("x"),
		NewNumber(42),
		NewList(types.NewUnknown(nil), nil),
		NewList(types.NewUnknown(nil), []ListArgument{NewListArgument(NewNumber(42), false)}),
		NewList(
			types.NewUnknown(nil),
			[]ListArgument{
				NewListArgument(NewNumber(42), false),
				NewListArgument(NewNumber(42), false),
			},
		),
		NewList(
			types.NewUnknown(nil),
			[]ListArgument{
				NewListArgument(
					NewList(types.NewUnknown(nil), []ListArgument{NewListArgument(NewNumber(42), false)}),
					false,
				),
			},
		),
		NewList(
			types.NewUnknown(nil),
			[]ListArgument{NewListArgument(NewVariable("xs"), true)},
		),
	} {
		assert.True(t, PatternsEqual(e, e))
	}

	for _, es := range [][2]Expression{
		{NewNumber(42), NewNumber(0)},
		{NewList(types.NewUnknown(nil), nil), NewNumber(42)},
		{
			NewList(types.NewUnknown(nil), []ListArgument{NewListArgument(NewNumber(42), false)}),
			NewList(types.NewUnknown(nil), []ListArgument{NewListArgument(NewNumber(0), false)}),
		},
		{
			NewList(types.NewUnknown(nil), []ListArgument{NewListArgument(NewNumber(42), false)}),
			NewList(
				types.NewUnknown(nil),
				[]ListArgument{
					NewListArgument(NewNumber(42), false),
					NewListArgument(NewNumber(42), false),
				},
			),
		},
		{
			NewList(
				types.NewUnknown(nil),
				[]ListArgument{NewListArgument(NewVariable("xs"), true)},
			),
			NewList(
				types.NewUnknown(nil),
				[]ListArgument{NewListArgument(NewVariable("xs"), false)},
			),
		},
	} {
		assert.False(t, PatternsEqual(es[0], es[1]))
	}
}

func TestPatternsEqualWithVariables(t *testing.T) {
	assert.True(t, PatternsEqual(NewVariable("x"), NewVariable("y")))
}
