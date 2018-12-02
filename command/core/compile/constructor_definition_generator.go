package compile

import (
	"github.com/raviqqe/jsonxx/command/core/ast"
	"github.com/raviqqe/jsonxx/command/core/compile/llir"
	"github.com/raviqqe/jsonxx/command/core/compile/names"
	"llvm.org/llvm/bindings/go/llvm"
)

type constructorDefinitionGenerator struct {
	module        llvm.Module
	typeGenerator typeGenerator
}

func newConstructorDefinitionGenerator(m llvm.Module) constructorDefinitionGenerator {
	return constructorDefinitionGenerator{m, newTypeGenerator(m)}
}

func (g constructorDefinitionGenerator) GenerateUnionifyFunction(d ast.ConstructorDefinition) error {
	f := llir.AddFunction(
		g.module,
		names.ToUnionify(d.Name()),
		g.typeGenerator.GenerateConstructorUnionifyFunction(d.Type(), d.Index()),
	)

	b := llvm.NewBuilder()
	b.SetInsertPointAtEnd(llvm.AddBasicBlock(f, ""))

	if len(d.Type().Constructors()) == 1 {
		b.CreateAggregateRet(f.Params())
	} else {
		p := b.CreateAlloca(f.Type().ElementType().ReturnType(), "")

		b.CreateStore(
			llvm.ConstInt(g.typeGenerator.GenerateConstructorTag(), uint64(d.Index()), false),
			b.CreateStructGEP(p, 0, ""),
		)

		pp := b.CreateBitCast(
			b.CreateStructGEP(p, 1, ""),
			llir.PointerType(
				g.typeGenerator.GenerateConstructorElements(d.Type().Constructors()[d.Index()]),
			),
			"",
		)

		for i, v := range f.Params() {
			b.CreateStore(v, b.CreateStructGEP(pp, i, ""))
		}

		b.CreateRet(b.CreateLoad(p, ""))
	}

	return llvm.VerifyFunction(f, llvm.AbortProcessAction)
}

func (g constructorDefinitionGenerator) GenerateStructifyFunction(d ast.ConstructorDefinition) error {
	f := llir.AddFunction(
		g.module,
		names.ToStructify(d.Name()),
		g.typeGenerator.GenerateConstructorStructifyFunction(d.Type(), d.Index()),
	)

	b := llvm.NewBuilder()
	b.SetInsertPointAtEnd(llvm.AddBasicBlock(f, ""))

	if len(d.Type().Constructors()) == 1 {
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

func (g constructorDefinitionGenerator) GenerateTag(d ast.ConstructorDefinition) {
	t := g.typeGenerator.GenerateConstructorTag()
	v := llvm.AddGlobal(g.module, t, names.ToTag(d.Name()))
	v.SetInitializer(llvm.ConstInt(t, uint64(d.Index()), false))
	v.SetGlobalConstant(true)
	v.SetUnnamedAddr(true)
}
