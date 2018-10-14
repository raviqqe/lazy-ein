package codegen

import (
	"testing"

	"github.com/raviqqe/stg/ast"
	"github.com/stretchr/testify/assert"
	"llvm.org/llvm/bindings/go/llvm"
)

func TestGenerate(t *testing.T) {
	m, err := Generate(
		"foo",
		[]ast.Bind{ast.NewBind("foo", ast.NewLambda(nil, true, nil, ast.NewFloat64(42)))},
	)

	assert.NotEqual(t, llvm.Module{}, m)
	assert.Nil(t, err)
}
