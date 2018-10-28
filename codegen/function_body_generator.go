package codegen

import (
	"fmt"

	"github.com/raviqqe/stg/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

const environmentArgumentName = "environment"

type functionBodyGenerator struct {
	function  llvm.Value
	builder   llvm.Builder
	variables map[string]llvm.Value
}

func newFunctionBodyGenerator(f llvm.Value, b llvm.Builder, as []string, gs map[string]llvm.Value) *functionBodyGenerator {
	vs := make(map[string]llvm.Value, len(as)+len(gs))

	for k, v := range gs {
		vs[k] = v
	}

	for i, n := range append([]string{environmentArgumentName}, as...) {
		v := f.Param(i)
		v.SetName(n)
		vs[n] = v
	}

	return &functionBodyGenerator{f, b, vs}
}

func (g *functionBodyGenerator) Generate(e ast.Expression) (llvm.Value, error) {
	g.builder.SetInsertPointAtEnd(llvm.AddBasicBlock(g.function, ""))
	return g.generateExpression(e)
}

func (g *functionBodyGenerator) generateExpression(e ast.Expression) (llvm.Value, error) {
	switch e := e.(type) {
	case ast.Application:
		f, err := g.resolveName(e.Function().Name())

		if err != nil {
			return llvm.Value{}, err
		} else if len(e.Arguments()) == 0 {
			return f, nil
		}

		vs := make([]llvm.Value, 0, len(e.Arguments()))

		for _, a := range e.Arguments() {
			switch a := a.(type) {
			case ast.Float64:
				vs = append(vs, llvm.ConstFloat(llvm.DoubleType(), a.Value()))
			default:
				vs = append(vs, g.variables[a.(ast.Variable).Name()])
			}
		}

		return g.builder.CreateCall(
			g.builder.CreateLoad(g.builder.CreateStructGEP(f, 0, ""), ""),
			append([]llvm.Value{g.builder.CreateStructGEP(f, 1, "")}, vs...),
			"",
		), nil
	case ast.Float64:
		return llvm.ConstFloat(llvm.DoubleType(), e.Value()), nil
	}

	panic("unreachable")
}

func (g *functionBodyGenerator) resolveName(n string) (llvm.Value, error) {
	v, ok := g.variables[n]

	if !ok {
		return llvm.Value{}, fmt.Errorf(`variable "%v" not found`, n)
	}

	return v, nil
}
