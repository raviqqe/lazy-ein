package codegen

import (
	"testing"

	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/llir"
	"github.com/raviqqe/stg/types"
	"github.com/stretchr/testify/assert"
	"llvm.org/llvm/bindings/go/llvm"
)

func TestFunctionBodyGeneratorGenerate(t *testing.T) {
	f := llvm.AddFunction(
		llvm.NewModule("foo"),
		"foo",
		llir.FunctionType(
			llvm.DoubleType(),
			[]llvm.Type{types.NewEnvironment(0).LLVMPointerType()},
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

func TestFunctionBodyGeneratorLambdaToEnvironment(t *testing.T) {
	for _, c := range []struct {
		lambda ast.Lambda
		size   int
	}{
		{
			lambda: ast.NewLambda(nil, true, nil, ast.NewFloat64(42), types.NewFloat64()),
			size:   8,
		},
		{
			lambda: ast.NewLambda(
				[]ast.Argument{
					ast.NewArgument("x", types.NewFloat64()),
					ast.NewArgument("y", types.NewFloat64()),
				},
				true,
				nil,
				ast.NewPrimitiveOperation(
					ast.AddFloat64,
					[]ast.Atom{ast.NewVariable("x"), ast.NewVariable("y")},
				),
				types.NewFloat64(),
			),
			size: 16,
		},
		{
			lambda: ast.NewLambda(nil, false, nil, ast.NewFloat64(42), types.NewFloat64()),
			size:   0,
		},
	} {
		f := llvm.AddFunction(
			llvm.NewModule("foo"),
			"foo",
			llir.FunctionType(
				llvm.DoubleType(),
				[]llvm.Type{types.NewEnvironment(0).LLVMPointerType()},
			),
		)

		b := llvm.NewBuilder()
		b.SetInsertPointAtEnd(llvm.AddBasicBlock(f, ""))

		e := newFunctionBodyGenerator(
			b,
			nil,
			nil,
		).lambdaToEnvironment(c.lambda)

		assert.Equal(t, c.size, e.LLVMType().ArrayLength())
	}
}
