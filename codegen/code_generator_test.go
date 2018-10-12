package codegen

import (
	"testing"

	"github.com/raviqqe/stg/ast"
)

func TestNewCodeGenerator(t *testing.T) {
	NewCodeGenerator("foo")
}

func TestCodeGeneratorGenerateModule(t *testing.T) {
	for _, bs := range [][]ast.Bind{
		nil,
		{ast.NewBind("foo", ast.NewLambda(nil, true, nil, ast.NewFloat64(42)))},
	} {
		NewCodeGenerator("foo").GenerateModule(bs)
	}
}
