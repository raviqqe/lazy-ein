package codegen

import (
	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/types"
	"llvm.org/llvm/bindings/go/llvm"
)

type functionGenerator struct {
	function    llvm.Value
	builder     llvm.Builder
	namedValues map[string]llvm.Value
	body        ast.Expression
}

func newFunctionGenerator(b ast.Bind, m llvm.Module) *functionGenerator {
	f := llvm.AddFunction(
		m,
		toEntryName(b.Name()),
		llvm.FunctionType(
			b.Lambda().ResultType().LLVMType(),
			append(
				[]llvm.Type{environmentPointerType},
				types.ToLLVMTypes(b.Lambda().ArgumentTypes())...,
			),
			false,
		),
	)

	vs := make(map[string]llvm.Value, len(b.Lambda().ArgumentNames()))

	for i, n := range append([]string{"environment"}, b.Lambda().ArgumentNames()...) {
		v := f.Param(i)
		v.SetName(n)
		vs[n] = v
	}

	return &functionGenerator{f, llvm.NewBuilder(), vs, b.Lambda().Body()}
}

func (g *functionGenerator) Generate() {
	b := llvm.AddBasicBlock(g.function, "")
	g.builder.SetInsertPointAtEnd(b)
	g.builder.CreateRet(g.generateExpression(g.body))
}

func (g *functionGenerator) generateExpression(e ast.Expression) llvm.Value {
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
			vs,
			"",
		)
	case ast.Float64:
		return llvm.ConstFloat(llvm.DoubleType(), e.Value())
	}

	panic("unreachable")
}
