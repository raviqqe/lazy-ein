package codegen

import (
	"errors"
	"fmt"

	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/types"
	"llvm.org/llvm/bindings/go/llvm"
)

type functionBodyGenerator struct {
	builder      llvm.Builder
	variables    map[string]llvm.Value
	createLambda func(string, ast.Lambda) (llvm.Value, error)
}

func newFunctionBodyGenerator(
	b llvm.Builder,
	vs map[string]llvm.Value,
	c func(string, ast.Lambda) (llvm.Value, error),
) *functionBodyGenerator {
	return &functionBodyGenerator{b, vs, c}
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
	case ast.Let:
		return g.generateLet(e)
	case ast.PrimitiveOperation:
		return g.generatePrimitiveOperation(e)
	}

	panic("unreachable")
}

func (g *functionBodyGenerator) generateLet(l ast.Let) (llvm.Value, error) {
	vs := g.saveVariables()
	defer g.restoreVariables(vs)

	for _, b := range l.Binds() {
		p := g.builder.CreateMalloc(
			llvm.StructType(
				[]llvm.Type{
					llvm.PointerType(
						llvm.FunctionType(
							types.Unbox(b.Lambda().ResultType()).LLVMType(),
							append(
								[]llvm.Type{types.NewEnvironment(0).LLVMPointerType()},
								types.ToLLVMTypes(b.Lambda().ArgumentTypes())...,
							),
							false,
						),
						0,
					),
					g.lambdaToEnvironment(b.Lambda()).LLVMType(),
				},
				false,
			),
			"",
		)

		f, err := g.createLambda(
			toInternalLambdaName(g.function().Name(), b.Name()),
			b.Lambda(),
		)

		if err != nil {
			return llvm.Value{}, err
		}

		g.builder.CreateStore(f, g.builder.CreateStructGEP(p, 0, ""))

		e := g.builder.CreateBitCast(
			g.builder.CreateStructGEP(p, 1, ""),
			llvm.PointerType(g.lambdaToFreeVariablesStructType(b.Lambda()), 0),
			"",
		)

		for i, n := range b.Lambda().FreeVariableNames() {
			v, err := g.resolveVariable(ast.Variable(n))

			if err != nil {
				return llvm.Value{}, err
			}

			g.builder.CreateStore(v, g.builder.CreateStructGEP(e, i, ""))
		}

		g.variables[b.Name()] = p
	}

	return g.generateExpression(l.Expression())
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

func (g *functionBodyGenerator) saveVariables() map[string]llvm.Value {
	vs := make(map[string]llvm.Value, len(g.variables))

	for k, v := range g.variables {
		vs[k] = v
	}

	return vs
}

func (g *functionBodyGenerator) restoreVariables(vs map[string]llvm.Value) {
	g.variables = vs
}

func (g *functionBodyGenerator) lambdaToEnvironment(l ast.Lambda) types.Environment {
	n := g.typeSize(g.lambdaToFreeVariablesStructType(l))

	if m := g.typeSize(types.Unbox(l.ResultType()).LLVMType()); m > n {
		n = m
	}

	return types.NewEnvironment(n)
}

func (functionBodyGenerator) lambdaToFreeVariablesStructType(l ast.Lambda) llvm.Type {
	return llvm.StructType(types.ToLLVMTypes(l.FreeVariableTypes()), false)
}

func (g *functionBodyGenerator) typeSize(t llvm.Type) int {
	return typeSize(g.function().GlobalParent(), t)
}

func (g *functionBodyGenerator) function() llvm.Value {
	return g.builder.GetInsertBlock().Parent()
}
