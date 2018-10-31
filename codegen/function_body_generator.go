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
		return g.generateApplication(e)
	case ast.Float64:
		return llvm.ConstFloat(llvm.DoubleType(), e.Value()), nil
	case ast.Let:
		return g.generateLet(e)
	case ast.PrimitiveOperation:
		return g.generatePrimitiveOperation(e)
	}

	panic("unreachable")
}

func (g *functionBodyGenerator) generateApplication(a ast.Application) (llvm.Value, error) {
	f, err := g.resolveName(a.Function().Name())

	if err != nil {
		return llvm.Value{}, err
	} else if len(a.Arguments()) == 0 {
		return f, nil
	}

	vs := make([]llvm.Value, 0, len(a.Arguments()))

	for _, a := range a.Arguments() {
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
}

func (g *functionBodyGenerator) generateLet(l ast.Let) (llvm.Value, error) {
	vs := make(map[string]llvm.Value, len(l.Binds()))

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
			llvm.PointerType(lambdaToFreeVariablesStructType(b.Lambda()), 0),
			"",
		)

		for i, n := range b.Lambda().FreeVariableNames() {
			v, err := g.resolveName(n)

			if err != nil {
				return llvm.Value{}, err
			}

			g.builder.CreateStore(v, g.builder.CreateStructGEP(e, i, ""))
		}

		vs[b.Name()] = p
	}

	return g.addVariables(vs).generateExpression(l.Expression())
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

func (g *functionBodyGenerator) resolveName(s string) (llvm.Value, error) {
	v, ok := g.variables[s]

	if !ok {
		return llvm.Value{}, fmt.Errorf(`variable "%v" not found`, s)
	}

	return v, nil
}

func (g *functionBodyGenerator) generateAtom(a ast.Atom) (llvm.Value, error) {
	switch a := a.(type) {
	case ast.Float64:
		return llvm.ConstFloat(llvm.DoubleType(), a.Value()), nil
	default:
		return g.resolveName(a.(ast.Variable).Name())
	}
}

func (g *functionBodyGenerator) generatePrimitiveArguments(as []ast.Atom) (llvm.Value, llvm.Value, error) {
	if len(as) != 2 {
		return llvm.Value{}, llvm.Value{}, errors.New(
			"invalid number of arguments to a binary primitive operation",
		)
	}

	v, err := g.generateAtom(as[0])

	if err != nil {
		return llvm.Value{}, llvm.Value{}, err
	}

	vv, err := g.generateAtom(as[1])

	if err != nil {
		return llvm.Value{}, llvm.Value{}, err
	}

	return v, vv, nil
}

func (g *functionBodyGenerator) addVariables(vs map[string]llvm.Value) *functionBodyGenerator {
	vvs := make(map[string]llvm.Value, len(g.variables)+len(vs))

	for k, v := range g.variables {
		vvs[k] = v
	}

	for k, v := range vs {
		vvs[k] = v
	}

	return &functionBodyGenerator{g.builder, vvs, g.createLambda}
}

func (g functionBodyGenerator) lambdaToEnvironment(l ast.Lambda) types.Environment {
	n := g.typeSize(lambdaToFreeVariablesStructType(l))

	if m := g.typeSize(types.Unbox(l.ResultType()).LLVMType()); m > n {
		n = m
	}

	return types.NewEnvironment(n)
}

func (g functionBodyGenerator) typeSize(t llvm.Type) int {
	return typeSize(g.function().GlobalParent(), t)
}

func (g functionBodyGenerator) function() llvm.Value {
	return g.builder.GetInsertBlock().Parent()
}
