package codegen

import (
	"github.com/raviqqe/stg/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

const environmentArgumentName = "environment"

type functionBodyGenerator struct {
	function    llvm.Value
	builder     llvm.Builder
	namedValues map[string]llvm.Value
}

func newFunctionBodyGenerator(f llvm.Value, b llvm.Builder, as []string) *functionBodyGenerator {
	vs := make(map[string]llvm.Value, len(as))

	for i, n := range append([]string{environmentArgumentName}, as...) {
		v := f.Param(i)
		v.SetName(n)
		vs[n] = v
	}

	return &functionBodyGenerator{f, b, vs}
}

func (g *functionBodyGenerator) Generate(e ast.Expression) llvm.Value {
	g.builder.SetInsertPointAtEnd(llvm.AddBasicBlock(g.function, ""))
	return g.generateExpression(e)
}

func (g *functionBodyGenerator) generateExpression(e ast.Expression) llvm.Value {
	switch e := e.(type) {
	case ast.Application:
		f := g.namedValues[e.Function().Name()]

		if len(e.Arguments()) == 0 {
			return f
		}

		vs := make([]llvm.Value, 0, len(e.Arguments()))

		for _, a := range e.Arguments() {
			switch a := a.(type) {
			case ast.Float64:
				vs = append(vs, llvm.ConstFloat(llvm.DoubleType(), a.Value()))
			default:
				vs = append(vs, g.namedValues[a.(ast.Variable).Name()])
			}
		}

		return g.builder.CreateCall(
			g.builder.CreateLoad(g.builder.CreateStructGEP(f, 0, ""), ""),
			append([]llvm.Value{g.builder.CreateStructGEP(f, 1, "")}, vs...),
			"",
		)
	case ast.Float64:
		return llvm.ConstFloat(llvm.DoubleType(), e.Value())
	}

	panic("unreachable")
}
