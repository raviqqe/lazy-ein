package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBox(t *testing.T) {
	a := NewAlgebraic(NewConstructor())

	for _, tt := range []Type{a, NewBoxed(a)} {
		assert.True(t, Equal(NewBoxed(a), Box(tt)))
	}
}

func TestUnbox(t *testing.T) {
	a := NewAlgebraic(NewConstructor())

	for _, tt := range []Type{a, NewBoxed(a)} {
		assert.True(t, Equal(a, Unbox(tt)))
	}
}

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
			NewFunction([]Type{NewFloat64()}, NewAlgebraic(NewConstructor(NewIndex(1)))),
			NewFunction(
				[]Type{NewFloat64()},
				NewAlgebraic(
					NewConstructor(
						NewFunction([]Type{NewFloat64()}, NewAlgebraic(NewConstructor(NewIndex(1)))),
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
			NewFunction([]Type{NewFloat64()}, NewAlgebraic(NewConstructor())),
		},
		{
			NewFunction([]Type{NewFloat64()}, NewFloat64()),
			NewFunction([]Type{NewFloat64(), NewFloat64()}, NewFloat64()),
		},
	} {
		assert.False(t, Equal(ts[0], ts[1]))
	}
}

func TestEqualWithEnvironment(t *testing.T) {
	assert.True(
		t,
		EqualWithEnvironment(
			NewAlgebraic(NewConstructor(NewIndex(0))),
			NewAlgebraic(NewConstructor(NewIndex(0))),
			[]Type{NewAlgebraic(NewConstructor(NewAlgebraic(NewConstructor(NewIndex(0)))))},
		),
	)
	assert.True(
		t,
		EqualWithEnvironment(
			NewAlgebraic(NewConstructor(NewIndex(1))),
			NewAlgebraic(NewConstructor(NewIndex(1))),
			[]Type{NewAlgebraic(NewConstructor(NewAlgebraic(NewConstructor(NewIndex(1)))))},
		),
	)
	assert.True(
		t,
		EqualWithEnvironment(
			NewAlgebraic(NewConstructor(NewIndex(0))),
			NewAlgebraic(NewConstructor(NewIndex(1))),
			[]Type{NewAlgebraic(NewConstructor(NewAlgebraic(NewConstructor(NewIndex(0)))))},
		),
	)

	assert.False(
		t,
		EqualWithEnvironment(
			NewAlgebraic(NewConstructor(NewIndex(0)), NewConstructor()),
			NewAlgebraic(NewConstructor(NewIndex(1))),
			[]Type{NewAlgebraic(NewConstructor(NewAlgebraic(NewConstructor(NewIndex(0)))))},
		),
	)
}

func TestValidate(t *testing.T) {
	for _, tt := range []Type{
		NewAlgebraic(NewConstructor()),
		NewAlgebraic(NewConstructor(NewIndex(0))),
		NewAlgebraic(NewConstructor(NewBoxed(NewIndex(0)))),
		NewFunction([]Type{NewFloat64()}, NewFloat64()),
		NewFunction([]Type{NewIndex(0)}, NewFloat64()),
	} {
		assert.True(t, Validate(tt))
	}

	for _, tt := range []Type{
		NewAlgebraic(NewConstructor(NewIndex(1))),
		NewFunction([]Type{NewIndex(1)}, NewFloat64()),
	} {
		assert.False(t, Validate(tt))
	}
}

func TestIsRecursive(t *testing.T) {
	for _, tt := range []Type{
		NewAlgebraic(NewConstructor(NewIndex(0))),
		NewAlgebraic(NewConstructor(NewBoxed(NewIndex(0)))),
		NewFunction([]Type{NewIndex(0)}, NewFloat64()),
	} {
		assert.True(t, IsRecursive(tt))
	}

	for _, tt := range []Type{
		NewAlgebraic(NewConstructor()),
		NewFunction([]Type{NewFloat64()}, NewFloat64()),
		NewAlgebraic(NewConstructor(NewAlgebraic(NewConstructor(NewIndex(0))))),
	} {
		assert.False(t, IsRecursive(tt))
	}
}

func TestUnwrap(t *testing.T) {
	for _, ts := range [][2]Type{
		{
			NewFloat64(),
			Unwrap(NewFloat64()),
		},
		{
			NewAlgebraic(NewConstructor()),
			Unwrap(NewAlgebraic(NewConstructor())),
		},
		{
			NewAlgebraic(NewConstructor(NewBoxed(NewIndex(0)))),
			Unwrap(
				NewAlgebraic(NewConstructor(NewBoxed(NewIndex(0)))),
			).(Algebraic).Constructors()[0].Elements()[0].(Boxed).Content(),
		},
		{
			NewFunction([]Type{NewIndex(0)}, NewFloat64()),
			Unwrap(NewFunction([]Type{NewIndex(0)}, NewFloat64())).(Function).Arguments()[0],
		},
		{
			NewBoxed(NewAlgebraic(NewConstructor(NewIndex(0)))),
			Unwrap(
				NewAlgebraic(NewConstructor(NewBoxed(NewAlgebraic(NewConstructor(NewIndex(0)))))),
			).(Algebraic).Constructors()[0].Elements()[0],
		},
	} {
		assert.Equal(t, ts[0], ts[1])
	}
}
