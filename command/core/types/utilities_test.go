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
		// Algebraic types
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
		// Functions
		{
			NewFunction([]Type{NewFloat64()}, NewFloat64()),
			NewFunction([]Type{NewFloat64()}, NewFloat64()),
		},
		{
			NewFunction([]Type{NewIndex(0)}, NewFloat64()),
			NewFunction([]Type{NewIndex(0)}, NewFloat64()),
		},
		{
			NewFunction([]Type{NewIndex(0)}, NewFloat64()),
			NewFunction([]Type{NewFunction([]Type{NewIndex(0)}, NewFloat64())}, NewFloat64()),
		},
		{
			NewFunction([]Type{NewFloat64()}, NewAlgebraic(NewConstructor(NewIndex(0)))),
			NewFunction([]Type{NewFloat64()}, NewAlgebraic(NewConstructor(NewIndex(0)))),
		},
		{
			NewFunction([]Type{NewFloat64()}, NewAlgebraic(NewConstructor(NewIndex(0)))),
			NewFunction(
				[]Type{NewFloat64()},
				NewAlgebraic(
					NewConstructor(
						NewFunction([]Type{NewFloat64()}, NewAlgebraic(NewConstructor(NewIndex(0)))),
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
		// Algebraic types
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
		// Functions
		{
			NewFunction([]Type{NewFloat64()}, NewFloat64()),
			NewFunction([]Type{NewAlgebraic(NewConstructor())}, NewFloat64()),
		},
		{
			NewFunction([]Type{NewFloat64()}, NewFloat64()),
			NewFunction([]Type{NewFloat64(), NewFloat64()}, NewFloat64()),
		},
	} {
		assert.False(t, Equal(ts[0], ts[1]))
	}
}
