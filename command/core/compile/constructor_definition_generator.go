package compile

import (
	"github.com/ein-lang/ein/command/core/compile/llir"
	"github.com/ein-lang/ein/command/core/compile/names"
	"github.com/ein-lang/ein/command/core/types"
	"llvm.org/llvm/bindings/go/llvm"
)

type constructorDefinitionGenerator struct {
	module        llvm.Module
	typeGenerator typeGenerator
}

func newConstructorDefinitionGenerator(
	m llvm.Module,
	g typeGenerator,
) constructorDefinitionGenerator {
	return constructorDefinitionGenerator{m, g}
}

func (g constructorDefinitionGenerator) Generate(a types.Algebraic) error {
	for i, c := range a.Constructors() {
		if err := g.generateUnionifyFunction(a, c, i); err != nil {
			return err
		}

		if err := g.generateStructifyFunction(a, c); err != nil {
			return err
		}

		g.GenerateTag(c, i)
	}

	return nil
}

func (g constructorDefinitionGenerator) generateUnionifyFunction(
	a types.Algebraic,
	c types.Constructor,
	i int,
) error {
	f := llir.AddFunction(
		g.module,
		names.ToUnionify(c.Name()),
		g.typeGenerator.GenerateConstructorUnionifyFunction(a, c),
	)

	b := llvm.NewBuilder()
	b.SetInsertPointAtEnd(llvm.AddBasicBlock(f, ""))

	if len(a.Constructors()) == 1 {
		b.CreateAggregateRet(f.Params())
	} else {
		p := b.CreateAlloca(f.Type().ElementType().ReturnType(), "")

		b.CreateStore(
			llvm.ConstInt(g.typeGenerator.GenerateConstructorTag(), uint64(i), false),
			b.CreateStructGEP(p, 0, ""),
		)

		pp := b.CreateBitCast(
			b.CreateStructGEP(p, 1, ""),
			llir.PointerType(g.typeGenerator.GenerateConstructorElements(c)),
			"",
		)

		for i, v := range f.Params() {
			b.CreateStore(v, b.CreateStructGEP(pp, i, ""))
		}

		b.CreateRet(b.CreateLoad(p, ""))
	}

	return llvm.VerifyFunction(f, llvm.AbortProcessAction)
}

func (g constructorDefinitionGenerator) generateStructifyFunction(
	a types.Algebraic,
	c types.Constructor,
) error {
	f := llir.AddFunction(
		g.module,
		names.ToStructify(c.Name()),
		g.typeGenerator.GenerateConstructorStructifyFunction(a, c),
	)

	b := llvm.NewBuilder()
	b.SetInsertPointAtEnd(llvm.AddBasicBlock(f, ""))

	if len(a.Constructors()) == 1 {
		b.CreateRet(f.FirstParam())
	} else {
		p := b.CreateAlloca(f.FirstParam().Type(), "")
		b.CreateStore(f.FirstParam(), p)

		p = b.CreateStructGEP(
			b.CreateBitCast(
				p,
				llir.PointerType(
					llir.StructType(
						[]llvm.Type{
							g.typeGenerator.GenerateConstructorTag(),
							f.Type().ElementType().ReturnType(),
						},
					),
				),
				"",
			),
			1,
			"",
		)

		b.CreateRet(b.CreateLoad(p, ""))
	}

	return llvm.VerifyFunction(f, llvm.AbortProcessAction)
}

func (g constructorDefinitionGenerator) GenerateTag(c types.Constructor, i int) {
	t := g.typeGenerator.GenerateConstructorTag()
	v := llvm.AddGlobal(g.module, t, names.ToTag(c.Name()))
	v.SetInitializer(llvm.ConstInt(t, uint64(i), false))
	v.SetGlobalConstant(true)
	v.SetUnnamedAddr(true)
}
