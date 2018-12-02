package compile

import (
	"testing"

	"github.com/raviqqe/jsonxx/command/core/ast"
	"github.com/raviqqe/jsonxx/command/core/compile/llir"
	"github.com/stretchr/testify/assert"
	"llvm.org/llvm/bindings/go/llvm"
)

func TestFunctionBodyGeneratorGenerate(t *testing.T) {
	m := llvm.NewModule("foo")
	f := llir.AddFunction(
		m,
		"foo",
		llir.FunctionType(
			llvm.DoubleType(),
			[]llvm.Type{llir.PointerType(newTypeGenerator(m).GenerateUnsizedPayload())},
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
