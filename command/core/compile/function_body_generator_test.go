package compile

import (
	"testing"

	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/compile/llir"
	"github.com/stretchr/testify/assert"
	"llvm.org/llvm/bindings/go/llvm"
)

func TestFunctionBodyGeneratorGenerate(t *testing.T) {
	m := llvm.NewModule("foo")

	g, err := newTypeGenerator(m, nil)
	assert.Nil(t, err)

	f := llir.AddFunction(
		m,
		"foo",
		llir.FunctionType(
			llvm.DoubleType(),
			[]llvm.Type{llir.PointerType(g.GenerateUnsizedPayload())},
		),
	)

	b := llvm.NewBuilder()
	b.SetInsertPointAtEnd(llvm.AddBasicBlock(f, ""))

	v, err := newFunctionBodyGenerator(b, nil, nil, g).Generate(ast.NewFloat64(42))

	assert.Nil(t, err)
	assert.True(t, v.IsConstant())
}
