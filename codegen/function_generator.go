package codegen

import (
	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/types"
	"llvm.org/llvm/bindings/go/llvm"
)

const environment = "environment"

type functionGenerator struct {
	function    llvm.Value
	builder     llvm.Builder
	namedValues map[string]llvm.Value
	name        string
	lambda      ast.Lambda
	module      llvm.Module
}

func newFunctionGenerator(b ast.Bind, m llvm.Module) *functionGenerator {
	f := llvm.AddFunction(
		m,
		toEntryName(b.Name()),
		llvm.FunctionType(
			b.Lambda().ResultType().LLVMType(),
			append(
				[]llvm.Type{types.EnvironmentPointerType},
				types.ToLLVMTypes(b.Lambda().ArgumentTypes())...,
			),
			false,
		),
	)

	vs := make(map[string]llvm.Value, len(b.Lambda().ArgumentNames()))

	for i, n := range append([]string{environment}, b.Lambda().ArgumentNames()...) {
		v := f.Param(i)
		v.SetName(n)
		vs[n] = v
	}

	return &functionGenerator{
		f,
		llvm.NewBuilder(),
		vs,
		b.Name(),
		b.Lambda(),
		m,
	}
}

func (g *functionGenerator) Generate() {
	b := llvm.AddBasicBlock(g.function, "")
	g.builder.SetInsertPointAtEnd(b)
	v := g.generateExpression(g.lambda.Body())

	if g.lambda.IsUpdatable() {
		g.builder.CreateStore(
			v,
			g.builder.CreateBitCast(
				g.namedValues[environment],
				llvm.PointerType(v.Type(), 0),
				"",
			),
		)

		g.builder.CreateStore(
			generateUpdatedEntryFunction(
				g.module,
				toUpdatedEntryName(g.name),
				g.lambda.ResultType().LLVMType(),
			),
			g.builder.CreateGEP(
				g.builder.CreateBitCast(
					g.namedValues[environment],
					llvm.PointerType(
						llvm.PointerType(
							llvm.FunctionType(
								g.lambda.ResultType().LLVMType(),
								[]llvm.Type{types.EnvironmentPointerType},
								false,
							),
							0,
						),
						0,
					),
					"",
				),
				[]llvm.Value{llvm.ConstIntFromString(llvm.Int32Type(), "-1", 10)},
				"",
			),
		)
	}

	g.builder.CreateRet(v)
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
			append([]llvm.Value{g.builder.CreateStructGEP(f, 1, "")}, vs...),
			"",
		)
	case ast.Float64:
		return llvm.ConstFloat(llvm.DoubleType(), e.Value())
	}

	panic("unreachable")
}

func generateUpdatedEntryFunction(m llvm.Module, s string, t llvm.Type) llvm.Value {
	f := llvm.AddFunction(
		m,
		s,
		llvm.FunctionType(
			t,
			[]llvm.Type{types.EnvironmentPointerType},
			false,
		),
	)

	bb := llvm.AddBasicBlock(f, "")
	b := llvm.NewBuilder()
	b.SetInsertPointAtEnd(bb)
	b.CreateRet(b.CreateLoad(b.CreateBitCast(f.FirstParam(), llvm.PointerType(t, 0), ""), ""))

	return f
}
