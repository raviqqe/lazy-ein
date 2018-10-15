package codegen

import (
	"github.com/raviqqe/stg/ast"
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
			b.Lambda().ResultType(),
			append([]llvm.Type{thunkPointerType}, b.Lambda().ArgumentTypes()...),
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
	b := llvm.AddBasicBlock(g.function, "entry")
	g.builder.SetInsertPointAtEnd(b)
	g.builder.CreateRet(g.generateExpression(g.body))
}

func (g *functionGenerator) generateExpression(e ast.Expression) llvm.Value {
	switch e := e.(type) {
	case ast.Application:
		if len(e.Arguments()) != 0 {
			panic("unimplemented")
		}

		return g.namedValues[e.Function().Name()]
	case ast.Float64:
		return llvm.ConstFloat(llvm.DoubleType(), e.Value())
	}

	panic("unreachable")
}
