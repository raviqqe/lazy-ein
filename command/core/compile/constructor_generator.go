package compile

import (
	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/compile/llir"
	"github.com/ein-lang/ein/command/core/compile/names"
	"github.com/ein-lang/ein/command/core/types"
	"llvm.org/llvm/bindings/go/llvm"
)

type constructorGenerator struct {
	module        llvm.Module
	typeGenerator typeGenerator
}

func newConstructorGenerator(
	m llvm.Module,
	g typeGenerator,
) constructorGenerator {
	return constructorGenerator{m, g}
}

func (g constructorGenerator) Generate(a types.Algebraic) error {
	for i := range a.Constructors() {
		if err := g.generateUnionifyFunction(ast.NewConstructor(a, i)); err != nil {
			return err
		}

		if err := g.generateStructifyFunction(ast.NewConstructor(a, i)); err != nil {
			return err
		}

		g.GenerateTag(ast.NewConstructor(a, i))
	}

	return nil
}

func (g constructorGenerator) generateUnionifyFunction(c ast.Constructor) error {
	f := llir.AddFunction(
		g.module,
		names.ToUnionify(c.ID()),
		g.typeGenerator.GenerateConstructorUnionifyFunction(c.AlgebraicType(), c.ConstructorType()),
	)
	f.SetLinkage(llvm.LinkOnceODRLinkage)

	b := llvm.NewBuilder()
	b.SetInsertPointAtEnd(llvm.AddBasicBlock(f, ""))

	if len(c.AlgebraicType().Constructors()) == 1 {
		b.CreateAggregateRet(f.Params())
	} else {
		p := b.CreateAlloca(f.Type().ElementType().ReturnType(), "")

		b.CreateStore(
			llvm.ConstInt(g.typeGenerator.GenerateConstructorTag(), uint64(c.Index()), false),
			b.CreateStructGEP(p, 0, ""),
		)

		pp := b.CreateBitCast(
			b.CreateStructGEP(p, 1, ""),
			llir.PointerType(g.typeGenerator.GenerateConstructorElements(c.ConstructorType())),
			"",
		)

		for i, v := range f.Params() {
			b.CreateStore(v, b.CreateStructGEP(pp, i, ""))
		}

		b.CreateRet(b.CreateLoad(p, ""))
	}

	return llvm.VerifyFunction(f, llvm.AbortProcessAction)
}

func (g constructorGenerator) generateStructifyFunction(
	c ast.Constructor,
) error {
	f := llir.AddFunction(
		g.module,
		names.ToStructify(c.ID()),
		g.typeGenerator.GenerateConstructorStructifyFunction(
			c.AlgebraicType(),
			c.ConstructorType(),
		),
	)
	f.SetLinkage(llvm.LinkOnceODRLinkage)

	b := llvm.NewBuilder()
	b.SetInsertPointAtEnd(llvm.AddBasicBlock(f, ""))

	if len(c.AlgebraicType().Constructors()) == 1 {
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

func (g constructorGenerator) GenerateTag(c ast.Constructor) {
	t := g.typeGenerator.GenerateConstructorTag()

	v := llvm.AddGlobal(g.module, t, names.ToTag(c.ID()))
	v.SetLinkage(llvm.LinkOnceODRLinkage)

	v.SetInitializer(llvm.ConstInt(t, uint64(c.Index()), false))
	v.SetGlobalConstant(true)
	v.SetUnnamedAddr(true)
}
