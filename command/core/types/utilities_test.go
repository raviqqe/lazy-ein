package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {
	for _, ts := range [][2]Type{
		{
			NewFloat64(),
			NewFloat64(),
		},
		{
			NewAlgebraic(NewConstructor()),
			NewAlgebraic(NewConstructor()),
		},
		{
			NewAlgebraic(NewConstructor(NewIndex(0))),
			NewAlgebraic(NewConstructor(NewIndex(0))),
		},
		{
			NewAlgebraic(NewConstructor(NewIndex(0))),
			NewAlgebraic(NewConstructor(NewAlgebraic(NewConstructor(NewIndex(0))))),
		},
		{
			NewAlgebraic(NewConstructor(NewIndex(0))),
			NewAlgebraic(NewConstructor(NewAlgebraic(NewConstructor(NewIndex(1))))),
		},
		{
			NewAlgebraic(
				NewConstructor(),
				NewConstructor(NewIndex(0)),
			),
			NewAlgebraic(
				NewConstructor(),
				NewConstructor(
					NewAlgebraic(
						NewConstructor(),
						NewConstructor(NewIndex(0)),
					),
				),
			),
		},
		{
			NewAlgebraic(
				NewConstructor(),
				NewConstructor(NewIndex(0)),
			),
			NewAlgebraic(
				NewConstructor(),
				NewConstructor(
					NewAlgebraic(
						NewConstructor(),
						NewConstructor(NewIndex(1)),
					),
				),
			),
		},
	} {
		assert.True(t, Equal(ts[0], ts[1]))
	}

	for _, ts := range [][2]Type{
		{
			NewFloat64(),
			NewAlgebraic(NewConstructor()),
		},
		{
			NewAlgebraic(NewConstructor()),
			NewAlgebraic(NewConstructor(), NewConstructor()),
		},
		{
			NewAlgebraic(NewConstructor()),
			NewAlgebraic(NewConstructor(NewFloat64())),
		},
		{
			NewAlgebraic(NewConstructor(NewFloat64())),
			NewAlgebraic(NewConstructor(NewAlgebraic(NewConstructor()))),
		},
	} {
		assert.False(t, Equal(ts[0], ts[1]))
	}
}
