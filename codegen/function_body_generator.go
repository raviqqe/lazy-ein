package codegen

import (
	"errors"
	"fmt"

	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/codegen/llir"
	"github.com/raviqqe/stg/codegen/names"
	"github.com/raviqqe/stg/types"
	"llvm.org/llvm/bindings/go/llvm"
)

type functionBodyGenerator struct {
	builder       llvm.Builder
	createLambda  func(string, ast.Lambda) (llvm.Value, error)
	nameGenerator *names.NameGenerator
	typeGenerator typeGenerator
	variables     map[string]llvm.Value
}

func newFunctionBodyGenerator(
	b llvm.Builder,
	vs map[string]llvm.Value,
	c func(string, ast.Lambda) (llvm.Value, error),
) *functionBodyGenerator {
	return &functionBodyGenerator{
		b,
		c,
		names.NewNameGenerator(b.GetInsertBlock().Parent().Name()),
		newTypeGenerator(b.GetInsertBlock().Parent().GlobalParent()),
		vs,
	}
}

func (g *functionBodyGenerator) Generate(e ast.Expression) (llvm.Value, error) {
	return g.generateExpression(e)
}

func (g *functionBodyGenerator) generateExpression(e ast.Expression) (llvm.Value, error) {
	switch e := e.(type) {
	case ast.Application:
		return g.generateApplication(e)
	case ast.Case:
		return g.generateCase(e)
	case ast.Let:
		return g.generateLet(e)
	case ast.Literal:
		return g.generateLiteral(e), nil
	case ast.PrimitiveOperation:
		return g.generatePrimitiveOperation(e)
	}

	panic("unreachable")
}

func (g *functionBodyGenerator) generateApplication(a ast.Application) (llvm.Value, error) {
	// TODO: Convert recursive function applications into thunks to allow tail-call elimination.
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

	return llir.CreateCall(
		g.builder,
		g.builder.CreateLoad(g.builder.CreateStructGEP(f, 0, ""), ""),
		append(
			[]llvm.Value{
				g.builder.CreateBitCast(
					g.builder.CreateStructGEP(f, 1, ""),
					types.NewPayload(0).LLVMPointerType(),
					"",
				),
			},
			vs...,
		),
	), nil
}

func (g *functionBodyGenerator) generateCase(c ast.Case) (llvm.Value, error) {
	e, err := g.generateExpression(c.Expression())

	if err != nil {
		return llvm.Value{}, err
	} else if e.Type().TypeKind() == llvm.PointerTypeKind {
		e = forceThunk(g.builder, e)
	}

	vs := make([]llvm.Value, 0, len(c.Alternatives())+1)
	bs := make([]llvm.BasicBlock, 0, len(c.Alternatives())+1)

	d := llvm.AddBasicBlock(g.function(), "default")
	p := llvm.AddBasicBlock(g.function(), "phi")
	s := g.builder.CreateSwitch(e, d, len(c.Alternatives()))

	for _, a := range c.Alternatives() {
		b := llvm.AddBasicBlock(g.function(), "")
		g.builder.SetInsertPointAtEnd(b)

		switch a := a.(type) {
		case ast.PrimitiveAlternative:
			s.AddCase(g.generateLiteral(a.Literal()), b)
		default:
			panic("unreachable")
		}

		v, err := g.generateExpression(a.Expression())

		if err != nil {
			return llvm.Value{}, err
		}

		vs = append(vs, v)
		bs = append(bs, b)

		g.builder.CreateBr(p)
	}

	t := vs[len(vs)-1].Type()

	g.builder.SetInsertPointAtEnd(d)

	v := llvm.ConstNull(t)

	if a, ok := c.DefaultAlternative(); ok {
		v, err = g.addVariables(
			map[string]llvm.Value{a.Variable(): e},
		).generateExpression(a.Expression())

		if err != nil {
			return llvm.Value{}, err
		}
	} else {
		g.builder.CreateUnreachable()
	}

	g.builder.CreateBr(p)

	vs = append(vs, v)
	bs = append(bs, d)

	g.builder.SetInsertPointAtEnd(p)
	v = g.builder.CreatePHI(t, "")
	v.AddIncoming(vs, bs)

	return v, nil
}

func (g *functionBodyGenerator) generateLet(l ast.Let) (llvm.Value, error) {
	vs := make(map[string]llvm.Value, len(l.Binds()))

	for _, b := range l.Binds() {
		as := b.Lambda().ArgumentTypes()
		r := b.Lambda().ResultType()

		if b.Lambda().IsThunk() {
			r = types.Unbox(r)
		}

		vs[b.Name()] = g.builder.CreateBitCast(
			g.builder.CreateMalloc(
				types.NewClosure(g.lambdaToPayload(b.Lambda()), as, r).LLVMType(),
				"",
			),
			llir.PointerType(types.NewClosure(types.NewPayload(0), as, r).LLVMType()),
			"",
		)
	}

	g = g.addVariables(vs)

	for _, b := range l.Binds() {
		f, err := g.createLambda(
			g.nameGenerator.Generate(b.Name()),
			b.Lambda(),
		)

		if err != nil {
			return llvm.Value{}, err
		}

		p := vs[b.Name()]

		g.builder.CreateStore(f, g.builder.CreateStructGEP(p, 0, ""))

		e := g.builder.CreateBitCast(
			g.builder.CreateStructGEP(p, 1, ""),
			llir.PointerType(g.typeGenerator.GenerateFreeVariables(b.Lambda())),
			"",
		)

		for i, n := range b.Lambda().FreeVariableNames() {
			v, err := g.resolveName(n)

			if err != nil {
				return llvm.Value{}, err
			}

			g.builder.CreateStore(v, g.builder.CreateStructGEP(e, i, ""))
		}
	}

	return g.generateExpression(l.Expression())
}

func (g *functionBodyGenerator) generateLiteral(l ast.Literal) llvm.Value {
	switch l := l.(type) {
	case ast.Float64:
		return llvm.ConstFloat(llvm.DoubleType(), l.Value())
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

	return &functionBodyGenerator{
		g.builder,
		g.createLambda,
		g.nameGenerator,
		g.typeGenerator,
		vvs,
	}
}

func (g functionBodyGenerator) lambdaToPayload(l ast.Lambda) types.Payload {
	n := g.typeSize(g.typeGenerator.GenerateFreeVariables(l))

	if m := g.typeSize(types.Unbox(l.ResultType()).LLVMType()); l.IsUpdatable() && m > n {
		n = m
	}

	return types.NewPayload(n)
}

func (g functionBodyGenerator) typeSize(t llvm.Type) int {
	return g.typeGenerator.GetSize(t)
}

func (g functionBodyGenerator) function() llvm.Value {
	return g.builder.GetInsertBlock().Parent()
}
