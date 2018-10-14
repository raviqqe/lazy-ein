package codegen

import (
	"testing"

	"github.com/raviqqe/stg/ast"
)

func TestNewModuleGenerator(t *testing.T) {
	newModuleGenerator("foo")
}

func TestModuleGeneratorGenerate(t *testing.T) {
	for _, bs := range [][]ast.Bind{
		nil,
		{ast.NewBind("foo", ast.NewLambda(nil, true, nil, ast.NewFloat64(42)))},
	} {
		newModuleGenerator("foo").Generate(bs)
	}
}
