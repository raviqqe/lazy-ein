package codegen

import (
	"testing"

	"github.com/raviqqe/stg/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

func TestNewModuleGenerator(t *testing.T) {
	newModuleGenerator("foo")
}

func TestModuleGeneratorGenerate(t *testing.T) {
	for _, bs := range [][]ast.Bind{
		nil,
		{
			ast.NewBind("foo", ast.NewLambda(nil, true, nil, ast.NewFloat64(42), llvm.DoubleType())),
		},
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{ast.NewArgument("x", llvm.DoubleType())},
					ast.NewApplication(ast.NewVariable("x"), nil),
					llvm.DoubleType(),
				),
			),
		},
	} {
		newModuleGenerator("foo").Generate(bs)
	}
}
