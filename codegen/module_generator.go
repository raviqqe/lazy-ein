package codegen

import (
	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/types"
	"llvm.org/llvm/bindings/go/llvm"
)

type moduleGenerator struct {
	module          llvm.Module
	globalVariables map[string]llvm.Value
}

func newModuleGenerator(m llvm.Module) *moduleGenerator {
	return &moduleGenerator{m, map[string]llvm.Value{}}
}

func (g *moduleGenerator) Generate(bs []ast.Bind) error {
	for _, b := range bs {
		e := types.NewEnvironment(g.getTypeSize(b.Lambda().ResultType().LLVMType()))
		f, err := g.createLambda(b.Name(), b.Lambda(), e)

		if err != nil {
			return err
		} else if err := llvm.VerifyFunction(f, llvm.AbortProcessAction); err != nil {
			return err
		}

		g.createClosure(b.Name(), f, e)
	}

	return llvm.VerifyModule(g.module, llvm.AbortProcessAction)
}

func (g *moduleGenerator) createLambda(n string, l ast.Lambda, e types.Environment) (llvm.Value, error) {
	f := llvm.AddFunction(
		g.module,
		toEntryName(n),
		llvm.FunctionType(
			l.ResultType().LLVMType(),
			append(
				[]llvm.Type{e.LLVMPointerType()},
				types.ToLLVMTypes(l.ArgumentTypes())...,
			),
			false,
		),
	)

	b := llvm.NewBuilder()
	v, err := newFunctionBodyGenerator(f, b, l.ArgumentNames(), g.globalVariables).Generate(l.Body())

	if err != nil {
		return llvm.Value{}, err
	} else if l.IsUpdatable() {
		b.CreateStore(v, b.CreateBitCast(f.FirstParam(), llvm.PointerType(v.Type(), 0), ""))
		b.CreateStore(
			g.createUpdatedEntryFunction(n, l.ResultType().LLVMType(), e),
			g.environmentToEntryFunctionPointer(b, f.FirstParam(), l.ResultType().LLVMType(), e),
		)
	}

	b.CreateRet(v)

	return f, nil
}

func (g *moduleGenerator) createClosure(n string, f llvm.Value, e types.Environment) {
	llvm.AddGlobal(
		g.module,
		llvm.StructType([]llvm.Type{f.Type(), e.LLVMType()}, false),
		n,
	).SetInitializer(llvm.ConstStruct([]llvm.Value{f, llvm.ConstNull(e.LLVMType())}, false))

	g.globalVariables[n] = g.module.NamedGlobal(n)
}

func (g *moduleGenerator) createUpdatedEntryFunction(n string, t llvm.Type, e types.Environment) llvm.Value {
	f := llvm.AddFunction(
		g.module,
		toUpdatedEntryName(n),
		llvm.FunctionType(t, []llvm.Type{e.LLVMPointerType()}, false),
	)

	b := llvm.NewBuilder()
	b.SetInsertPointAtEnd(llvm.AddBasicBlock(f, ""))
	b.CreateRet(b.CreateLoad(b.CreateBitCast(f.FirstParam(), llvm.PointerType(t, 0), ""), ""))

	return f
}

func (g *moduleGenerator) getTypeSize(t llvm.Type) int {
	return int(llvm.NewTargetData(g.module.DataLayout()).TypeAllocSize(t))
}

func (g *moduleGenerator) environmentToEntryFunctionPointer(
	b llvm.Builder, v llvm.Value, t llvm.Type, e types.Environment,
) llvm.Value {
	return b.CreateGEP(
		b.CreateBitCast(
			v,
			llvm.PointerType(
				llvm.PointerType(
					llvm.FunctionType(t, []llvm.Type{e.LLVMPointerType()}, false),
					0,
				),
				0,
			),
			"",
		),
		[]llvm.Value{llvm.ConstIntFromString(llvm.Int32Type(), "-1", 10)},
		"",
	)
}
