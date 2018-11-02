package codegen

import (
	"testing"

	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/types"
	"github.com/stretchr/testify/assert"
	"llvm.org/llvm/bindings/go/llvm"
)

func TestFunctionBodyGeneratorGenerate(t *testing.T) {
	f := llvm.AddFunction(
		llvm.NewModule("foo"),
		"foo",
		llvm.FunctionType(
			llvm.DoubleType(),
			[]llvm.Type{types.NewEnvironment(0).LLVMPointerType()},
			false,
		),
	)

	b := llvm.NewBuilder()
	b.SetInsertPointAtEnd(llvm.AddBasicBlock(f, ""))

	v, err := newFunctionBodyGenerator(
		b,
		nil,
		nil,
	).Generate(ast.NewFloat64(42))

	assert.Nil(t, err)
	assert.True(t, v.IsConstant())
}
