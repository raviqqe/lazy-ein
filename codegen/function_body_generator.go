package codegen

import (
	"errors"
	"fmt"

	"github.com/raviqqe/stg/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

type functionBodyGenerator struct {
	builder   llvm.Builder
	variables map[string]llvm.Value
}

func newFunctionBodyGenerator(b llvm.Builder, vs map[string]llvm.Value) *functionBodyGenerator {
	return &functionBodyGenerator{b, vs}
}

func (g *functionBodyGenerator) Generate(e ast.Expression) (llvm.Value, error) {
	return g.generateExpression(e)
}

func (g *functionBodyGenerator) generateExpression(e ast.Expression) (llvm.Value, error) {
	switch e := e.(type) {
	case ast.Application:
		f, err := g.resolveVariable(e.Function())

		if err != nil {
			return llvm.Value{}, err
		} else if len(e.Arguments()) == 0 {
			return f, nil
		}

		vs := make([]llvm.Value, 0, len(e.Arguments()))

		for _, a := range e.Arguments() {
			v, err := g.generateAtom(a)

			if err != nil {
				return llvm.Value{}, err
			}

			vs = append(vs, v)
		}

		return g.builder.CreateCall(
			g.builder.CreateLoad(g.builder.CreateStructGEP(f, 0, ""), ""),
			append([]llvm.Value{g.builder.CreateStructGEP(f, 1, "")}, vs...),
			"",
		), nil
	case ast.Float64:
		return llvm.ConstFloat(llvm.DoubleType(), e.Value()), nil
	case ast.PrimitiveOperation:
		return g.generatePrimitiveOperation(e)
	}

	panic("unreachable")
}

func (g *functionBodyGenerator) generatePrimitiveOperation(o ast.PrimitiveOperation) (llvm.Value, error) {
	l, r, err := g.generatePrimitiveArguments(o.Arguments())

	if err != nil {
		return llvm.Value{}, err
	}

	switch o.Primitive() {
	case ast.AddFloat64:
		return g.builder.CreateFAdd(l, r, ""), nil
	case ast.SubtractFloat64:
		return g.builder.CreateFSub(l, r, ""), nil
	case ast.MultiplyFloat64:
		return g.builder.CreateFMul(l, r, ""), nil
	case ast.DivideFloat64:
		return g.builder.CreateFDiv(l, r, ""), nil
	}

	panic("unreachable")
}

func (g *functionBodyGenerator) resolveVariable(x ast.Variable) (llvm.Value, error) {
	v, ok := g.variables[x.Name()]

	if !ok {
		return llvm.Value{}, fmt.Errorf(`variable "%v" not found`, x.Name())
	}

	return v, nil
}

func (g *functionBodyGenerator) generateAtom(a ast.Atom) (llvm.Value, error) {
	switch a := a.(type) {
	case ast.Float64:
		return llvm.ConstFloat(llvm.DoubleType(), a.Value()), nil
	default:
		return g.resolveVariable(a.(ast.Variable))
	}
}

func (g *functionBodyGenerator) generatePrimitiveArguments(as []ast.Atom) (llvm.Value, llvm.Value, error) {
	if len(as) != 2 {
		return llvm.Value{}, llvm.Value{}, errors.New(
			"invalid number of arguments to a binary primitive operation",
		)
	}

	vs := make([]llvm.Value, 0, len(as))

	for _, a := range as {
		v, err := g.generateAtom(a)

		if err != nil {
			return llvm.Value{}, llvm.Value{}, err
		}

		vs = append(vs, v)
	}

	return vs[0], vs[1], nil
}
