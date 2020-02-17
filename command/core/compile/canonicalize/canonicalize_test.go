package canonicalize_test

import (
	"testing"

	"github.com/raviqqe/lazy-ein/command/core/ast"
	"github.com/raviqqe/lazy-ein/command/core/compile/canonicalize"
	"github.com/raviqqe/lazy-ein/command/core/types"
	"github.com/stretchr/testify/assert"
)

func TestCanonicalizeWithAlgebraicTypes(t *testing.T) {
	for _, ts := range [][2]types.Bindable{
		{
			types.NewAlgebraic(types.NewConstructor()),
			types.NewAlgebraic(types.NewConstructor()),
		},
		{
			types.NewAlgebraic(types.NewConstructor(types.NewFloat64())),
			types.NewAlgebraic(types.NewConstructor(types.NewFloat64())),
		},
		{
			types.NewAlgebraic(types.NewConstructor(types.NewIndex(0))),
			types.NewAlgebraic(types.NewConstructor(types.NewIndex(0))),
		},
		{
			types.NewAlgebraic(types.NewConstructor(types.NewBoxed(types.NewIndex(0)))),
			types.NewAlgebraic(types.NewConstructor(types.NewBoxed(types.NewIndex(0)))),
		},
		{
			types.NewAlgebraic(
				types.NewConstructor(
					types.NewAlgebraic(types.NewConstructor(types.NewIndex(0))),
				),
			),
			types.NewAlgebraic(types.NewConstructor(types.NewIndex(0))),
		},
		{
			types.NewAlgebraic(
				types.NewConstructor(
					types.NewAlgebraic(types.NewConstructor(types.NewIndex(1))),
				),
			),
			types.NewAlgebraic(types.NewConstructor(types.NewIndex(0))),
		},
		{
			types.NewAlgebraic(
				types.NewConstructor(
					types.NewAlgebraic(
						types.NewConstructor(types.NewIndex(0)),
						types.NewConstructor(types.NewIndex(0)),
					),
				),
				types.NewConstructor(types.NewIndex(0)),
			),
			types.NewAlgebraic(
				types.NewConstructor(types.NewIndex(0)),
				types.NewConstructor(types.NewIndex(0)),
			),
		},
		{
			types.NewAlgebraic(
				types.NewConstructor(
					types.NewBoxed(
						types.NewAlgebraic(types.NewConstructor(types.NewBoxed(types.NewIndex(0)))),
					),
				),
			),
			types.NewAlgebraic(types.NewConstructor(types.NewBoxed(types.NewIndex(0)))),
		},
		{
			types.NewAlgebraic(
				types.NewConstructor(
					types.NewAlgebraic(
						types.NewConstructor(types.NewAlgebraic(types.NewConstructor(types.NewIndex(0)))),
					),
				),
				types.NewConstructor(),
			),
			types.NewAlgebraic(
				types.NewConstructor(
					types.NewAlgebraic(types.NewConstructor(types.NewIndex(0))),
				),
				types.NewConstructor(),
			),
		},
	} {
		assert.Equal(
			t,
			ts[1],
			canonicalize.Canonicalize(
				ast.NewModule(
					nil,
					[]ast.Bind{
						ast.NewBind("x", ast.NewVariableLambda(nil, ast.NewFloat64(42), ts[0])),
					},
				),
			).Binds()[0].Lambda().ResultType(),
		)
	}
}

func TestCanonicalizeWithFunctions(t *testing.T) {
	for _, ts := range [][2]types.Function{
		{
			types.NewFunction([]types.Type{types.NewFloat64()}, types.NewFloat64()),
			types.NewFunction([]types.Type{types.NewFloat64()}, types.NewFloat64()),
		},
		{
			types.NewFunction([]types.Type{types.NewIndex(0)}, types.NewFloat64()),
			types.NewFunction([]types.Type{types.NewIndex(0)}, types.NewFloat64()),
		},
		{
			types.NewFunction(
				[]types.Type{types.NewFunction([]types.Type{types.NewIndex(0)}, types.NewFloat64())},
				types.NewFloat64(),
			),
			types.NewFunction([]types.Type{types.NewIndex(0)}, types.NewFloat64()),
		},
		{
			types.NewFunction(
				[]types.Type{types.NewFunction([]types.Type{types.NewIndex(1)}, types.NewFloat64())},
				types.NewFloat64(),
			),
			types.NewFunction([]types.Type{types.NewIndex(0)}, types.NewFloat64()),
		},
	} {
		assert.Equal(
			t,
			ts[1],
			canonicalize.Canonicalize(
				ast.NewModule(
					nil,
					[]ast.Bind{
						ast.NewBind(
							"f",
							ast.NewFunctionLambda(
								nil,
								[]ast.Argument{ast.NewArgument("x", ts[0])}, ast.NewFloat64(42),
								types.NewFloat64(),
							),
						),
					},
				),
			).Binds()[0].Lambda().ArgumentTypes()[0],
		)
	}
}
