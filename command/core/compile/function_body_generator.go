package compile

import (
	"errors"
	"fmt"

	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/compile/llir"
	"github.com/ein-lang/ein/command/core/compile/names"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
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
	g typeGenerator,
) *functionBodyGenerator {
	return &functionBodyGenerator{
		b,
		c,
		names.NewNameGenerator(b.GetInsertBlock().Parent().Name()),
		g,
		vs,
	}
}

func (g *functionBodyGenerator) Generate(e ast.Expression) (llvm.Value, error) {
	return g.generateExpression(e)
}

func (g *functionBodyGenerator) generateExpression(e ast.Expression) (llvm.Value, error) {
	switch e := e.(type) {
	case ast.FunctionApplication:
		return g.generateFunctionApplication(e)
	case ast.Case:
		return g.generateCase(e)
	case ast.ConstructorApplication:
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

func (g *functionBodyGenerator) generateFunctionApplication(
	a ast.FunctionApplication,
) (llvm.Value, error) {
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
	switch c := c.(type) {
	case ast.AlgebraicCase:
		return g.generateAlgebraicCase(c)
	case ast.PrimitiveCase:
		return g.generateFloatCase(c)
	}

	panic("unreachable")
}

func (g *functionBodyGenerator) generateAlgebraicCase(c ast.AlgebraicCase) (llvm.Value, error) {
	arg, err := g.generateExpression(c.Argument())

	if err != nil {
		return llvm.Value{}, err
	}

	arg = g.builder.CreateLoad(forceThunk(g.builder, arg, g.typeGenerator), "")
	tag := llvm.ConstInt(g.typeGenerator.GenerateConstructorTag(), 0, false)

	if len(c.Alternatives()) != 0 && len(c.Alternatives()[0].Constructor().AlgebraicType().Constructors()) != 1 {
		tag = g.builder.CreateExtractValue(arg, 0, "")
	}

	p := newPhiGenerator(llvm.AddBasicBlock(g.function(), "phi"))
	d := llvm.AddBasicBlock(g.function(), "default")
	s := g.builder.CreateSwitch(tag, d, len(c.Alternatives()))

	for i, a := range c.Alternatives() {
		b := llvm.AddBasicBlock(g.function(), fmt.Sprintf("case.%v", i))
		s.AddCase(g.module().NamedGlobal(names.ToTag(a.Constructor().ID())).Initializer(), b)
		g.builder.SetInsertPointAtEnd(b)

		es := llir.CreateCall(
			g.builder,
			g.module().NamedFunction(names.ToStructify(a.Constructor().ID())),
			[]llvm.Value{arg},
		)

		vs := make(map[string]llvm.Value, len(a.ElementNames()))

		for i, n := range a.ElementNames() {
			vs[n] = g.builder.CreateExtractValue(es, i, "")
		}

		v, err := g.addVariables(vs).generateExpression(a.Expression())

		if err != nil {
			return llvm.Value{}, err
		}

		p.CreateBr(g.builder, v)
	}

	// Pass down forced unboxed case arguments.
	if err = g.generateDefaultAlternative(c, arg, d, p); err != nil {
		return llvm.Value{}, err
	}

	return p.Generate(g.builder), nil
}

func (g *functionBodyGenerator) generateFloatCase(c ast.PrimitiveCase) (llvm.Value, error) {
	v, err := g.generateExpression(c.Argument())

	if err != nil {
		return llvm.Value{}, nil
	}

	b := g.builder.GetInsertBlock()
	p := newPhiGenerator(llvm.AddBasicBlock(g.function(), "phi"))

	for i, a := range c.Alternatives() {
		g.builder.SetInsertPointAtEnd(b)
		b = llvm.AddBasicBlock(g.function(), fmt.Sprintf("else.%v", i))
		bb := llvm.AddBasicBlock(g.function(), fmt.Sprintf("then.%v", i))

		g.builder.CreateCondBr(
			g.builder.CreateFCmp(llvm.FloatOEQ, v, g.generateLiteral(a.Literal()), ""),
			bb,
			b,
		)

		g.builder.SetInsertPointAtEnd(bb)
		v, err := g.generateExpression(a.Expression())

		if err != nil {
			return llvm.Value{}, err
		}

		p.CreateBr(g.builder, v)
	}

	if err = g.generateDefaultAlternative(c, v, b, p); err != nil {
		return llvm.Value{}, err
	}

	return p.Generate(g.builder), nil
}

func (g *functionBodyGenerator) generateDefaultAlternative(
	c ast.Case,
	v llvm.Value,
	b llvm.BasicBlock,
	p *phiGenerator,
) error {
	g.builder.SetInsertPointAtEnd(b)

	a, ok := c.DefaultAlternative()

	if !ok {
		g.builder.CreateCall(g.module().NamedFunction(panicFunctionName), nil, "")
		g.builder.CreateUnreachable()
		return nil
	}

	v, err := g.addVariables(
		map[string]llvm.Value{a.Variable(): v},
	).generateExpression(a.Expression())

	if err != nil {
		return err
	}

	p.CreateBr(g.builder, v)

	return nil
}

func (g *functionBodyGenerator) generateConstructor(c ast.ConstructorApplication) (llvm.Value, error) {
	vs, err := g.generateAtoms(c.Arguments())

	if err != nil {
		return llvm.Value{}, err
	}

	v := llir.CreateCall(
		g.builder,
		g.module().NamedFunction(names.ToUnionify(c.Constructor().ID())),
		vs,
	)

	p := g.builder.CreateAlloca(v.Type(), "")
	g.builder.CreateStore(v, p)
	return p, nil
}

func (g *functionBodyGenerator) generateLet(l ast.Let) (llvm.Value, error) {
	vs := make(map[string]llvm.Value, len(l.Binds()))

	for _, b := range l.Binds() {
		t := g.typeGenerator.GenerateSizedClosure(b.Lambda().ToDeclaration())

		vs[b.Name()] = g.builder.CreateBitCast(
			g.allocateHeap(t),
			llir.PointerType(g.typeGenerator.GenerateUnsizedClosureFromSized(t)),
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
			llir.PointerType(g.typeGenerator.GenerateEnvironment(b.Lambda().ToDeclaration())),
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

	switch o.PrimitiveOperator() {
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

func (g *functionBodyGenerator) allocateHeap(t llvm.Type) llvm.Value {
	return g.builder.CreateBitCast(
		g.builder.CreateCall(
			g.module().NamedFunction(allocFunctionName),
			[]llvm.Value{
				llvm.ConstInt(llir.WordType(), uint64(g.typeGenerator.GetSize(t)), false),
			},
			"",
		),
		llir.PointerType(t),
		"",
	)
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

func (g functionBodyGenerator) module() llvm.Module {
	return g.function().GlobalParent()
}
