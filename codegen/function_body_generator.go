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
	case ast.Constructor:
		return g.generateConstructor(e)
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

	vs, err := g.generateAtoms(a.Arguments())

	if err != nil {
		return llvm.Value{}, err
	}

	return llir.CreateCall(
		g.builder,
		g.builder.CreateLoad(g.builder.CreateStructGEP(f, 0, ""), ""),
		append(
			[]llvm.Value{
				g.builder.CreateBitCast(
					g.builder.CreateStructGEP(f, 1, ""),
					llir.PointerType(g.typeGenerator.GenerateUnsizedPayload()),
					"",
				),
			},
			vs...,
		),
	), nil
}

func (g *functionBodyGenerator) generateCase(c ast.Case) (llvm.Value, error) {
	v, e, err := g.generateSwitcheeAndExpression(c)

	if err != nil {
		return llvm.Value{}, nil
	}

	p := newPhiGenerator(llvm.AddBasicBlock(g.function(), "phi"), c)

	d := llvm.AddBasicBlock(g.function(), "default")
	s := g.builder.CreateSwitch(v, d, len(c.Alternatives()))

	for _, a := range c.Alternatives() {
		b := llvm.AddBasicBlock(g.function(), "")
		g.builder.SetInsertPointAtEnd(b)

		v, g, err := g.generateAlternative(e, a)

		if err != nil {
			return llvm.Value{}, err
		}

		s.AddCase(v, b)

		v, err = g.generateExpression(a.Expression())

		if err != nil {
			return llvm.Value{}, err
		}

		g.builder.CreateBr(p.Block())
		p.AddIncoming(v, b)
	}

	g.builder.SetInsertPointAtEnd(d)

	if a, ok := c.DefaultAlternative(); ok {
		v, err := g.addVariables(
			map[string]llvm.Value{a.Variable(): e},
		).generateExpression(a.Expression())

		if err != nil {
			return llvm.Value{}, err
		}

		g.builder.CreateBr(p.Block())
		p.AddIncoming(v, d)
	} else {
		g.builder.CreateUnreachable()
	}

	return p.Generate(g.builder), nil
}

func (g *functionBodyGenerator) generateSwitcheeAndExpression(c ast.Case) (llvm.Value, llvm.Value, error) {
	v, err := g.generateExpression(c.Expression())

	if err != nil {
		return llvm.Value{}, llvm.Value{}, err
	} else if _, ok := c.ExpressionType().(types.Boxed); ok {
		v = forceThunk(g.builder, v, g.typeGenerator)
	}

	t, ok := c.ExpressionType().(types.Algebraic)

	if !ok {
		return v, v, nil
	} else if len(t.Constructors()) == 1 {
		return llvm.ConstInt(g.typeGenerator.GenerateConstructorTag(), 0, false), v, nil
	}

	return g.builder.CreateExtractValue(v, 0, ""), v, nil
}

func (g *functionBodyGenerator) generateAlternative(e llvm.Value, a ast.Alternative) (llvm.Value, *functionBodyGenerator, error) {
	switch a := a.(type) {
	case ast.PrimitiveAlternative:
		return g.generateLiteral(a.Literal()), g, nil
	case ast.AlgebraicAlternative:
		v, err := g.resolveName(names.ToTag(a.ConstructorName()))

		if err != nil {
			return llvm.Value{}, nil, err
		}

		f, err := g.resolveName(names.ToStructify(a.ConstructorName()))

		if err != nil {
			return llvm.Value{}, nil, err
		}

		es := g.builder.CreateCall(f, []llvm.Value{e}, "")
		vs := map[string]llvm.Value{}

		for i, n := range a.ElementNames() {
			vs[n] = g.builder.CreateExtractValue(es, i, "")
		}

		return v, g.addVariables(vs), nil
	}

	panic("unreachable")
}

func (g *functionBodyGenerator) generateConstructor(c ast.Constructor) (llvm.Value, error) {
	v, err := g.resolveName(names.ToUnionify(c.Name()))

	if err != nil {
		return llvm.Value{}, err
	}

	vs, err := g.generateAtoms(c.Arguments())

	if err != nil {
		return llvm.Value{}, err
	}

	return llir.CreateCall(g.builder, v, vs), nil
}

func (g *functionBodyGenerator) generateLet(l ast.Let) (llvm.Value, error) {
	vs := make(map[string]llvm.Value, len(l.Binds()))

	for _, b := range l.Binds() {
		vs[b.Name()] = g.builder.CreateBitCast(
			g.builder.CreateMalloc(g.typeGenerator.GenerateSizedClosure(b.Lambda()), ""),
			llir.PointerType(g.typeGenerator.GenerateUnsizedClosure(b.Lambda())),
			"",
		)
	}

	g = g.addVariables(vs)

	for _, b := range l.Binds() {
		f, err := g.createLambda(g.nameGenerator.Generate(b.Name()), b.Lambda())

		if err != nil {
			return llvm.Value{}, err
		}

		p := vs[b.Name()]

		g.builder.CreateStore(f, g.builder.CreateStructGEP(p, 0, ""))

		e := g.builder.CreateBitCast(
			g.builder.CreateStructGEP(p, 1, ""),
			llir.PointerType(g.typeGenerator.GenerateEnvironment(b.Lambda())),
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
	vs, err := g.generateAtoms(o.Arguments())

	if err != nil {
		return llvm.Value{}, err
	} else if len(vs) != 2 {
		return llvm.Value{}, errors.New("invalid number of arguments to a binary primitive operation")
	}

	switch o.Primitive() {
	case ast.AddFloat64:
		return g.builder.CreateFAdd(vs[0], vs[1], ""), nil
	case ast.SubtractFloat64:
		return g.builder.CreateFSub(vs[0], vs[1], ""), nil
	case ast.MultiplyFloat64:
		return g.builder.CreateFMul(vs[0], vs[1], ""), nil
	case ast.DivideFloat64:
		return g.builder.CreateFDiv(vs[0], vs[1], ""), nil
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

func (g *functionBodyGenerator) generateAtoms(as []ast.Atom) ([]llvm.Value, error) {
	vs := make([]llvm.Value, 0, len(as))

	for _, a := range as {
		v, err := g.generateAtom(a)

		if err != nil {
			return nil, err
		}

		vs = append(vs, v)
	}

	return vs, nil
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

func (g functionBodyGenerator) function() llvm.Value {
	return g.builder.GetInsertBlock().Parent()
}
