package codegen

import (
	"github.com/raviqqe/stg/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

type functionGenerator struct {
	function llvm.Value
	builder  llvm.Builder
}

func newFunctionGenerator(f llvm.Value) *functionGenerator {
	return &functionGenerator{f, llvm.NewBuilder()}
}

func (g *functionGenerator) Generate(e ast.Expression) {
	b := llvm.AddBasicBlock(g.function, "entry")
	g.builder.SetInsertPointAtEnd(b)
	g.builder.CreateRet(g.generateExpression(e))
}

func (g *functionGenerator) generateExpression(e ast.Expression) llvm.Value {
	switch e := e.(type) {
	case ast.Float64:
		return llvm.ConstFloat(llvm.DoubleType(), e.Value())
	}

	panic("unreachable")
}
