package codegen

import (
	"testing"

	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/types"
	"llvm.org/llvm/bindings/go/llvm"
)

func TestNewModuleGenerator(t *testing.T) {
	newModuleGenerator(llvm.NewModule("foo"))
}

func TestModuleGeneratorGenerate(t *testing.T) {
	for _, bs := range [][]ast.Bind{
		nil,
		{
			ast.NewBind("foo", ast.NewLambda(nil, true, nil, ast.NewFloat64(42), types.NewFloat64())),
		},
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewApplication(ast.NewVariable("x"), nil),
					types.NewFloat64(),
				),
			),
		},
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{ast.NewArgument("x", types.NewBoxed(types.NewFloat64()))},
					ast.NewApplication(ast.NewVariable("x"), nil),
					types.NewBoxed(types.NewFloat64()),
				),
			),
		},
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{
						ast.NewArgument(
							"x",
							types.NewFunction([]types.Type{types.NewFloat64()}, types.NewFloat64()),
						),
					},
					ast.NewApplication(ast.NewVariable("x"), []ast.Atom{ast.Float64(42)}),
					types.NewFloat64(),
				),
			),
		},
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{
						ast.NewArgument(
							"f",
							types.NewFunction([]types.Type{types.NewFloat64()}, types.NewFloat64()),
						),
						ast.NewArgument(
							"x",
							types.NewFloat64(),
						),
					},
					ast.NewApplication(ast.NewVariable("f"), []ast.Atom{ast.NewVariable("x")}),
					types.NewFloat64(),
				),
			),
		},
	} {
		newModuleGenerator(llvm.NewModule("foo")).Generate(bs)
	}
}
