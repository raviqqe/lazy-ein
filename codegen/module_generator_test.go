package codegen

import (
	"testing"

	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/types"
)

func TestNewModuleGenerator(t *testing.T) {
	newModuleGenerator("foo")
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
	} {
		newModuleGenerator("foo").Generate(bs)
	}
}
