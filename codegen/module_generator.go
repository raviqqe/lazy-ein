package codegen

import (
	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/types"
	"llvm.org/llvm/bindings/go/llvm"
)

type moduleGenerator struct {
	module llvm.Module
}

func newModuleGenerator(s string) *moduleGenerator {
	return &moduleGenerator{llvm.NewModule(s)}
}

func (g *moduleGenerator) Generate(bs []ast.Bind) error {
	for _, b := range bs {
		e := types.NewEnvironment(g.getTypeSize(b.Lambda().ResultType().LLVMType()))
		f := generateFunction(b, e, g.module)

		if err := llvm.VerifyFunction(f, llvm.AbortProcessAction); err != nil {
			return err
		}

		g.createClosure(b.Name(), f, e)
	}

	return llvm.VerifyModule(g.module, llvm.AbortProcessAction)
}

func (g *moduleGenerator) createClosure(s string, f llvm.Value, e types.Environment) {
	llvm.AddGlobal(
		g.module,
		llvm.StructType([]llvm.Type{f.Type(), e.LLVMType()}, false),
		s,
	).SetInitializer(llvm.ConstStruct([]llvm.Value{f, llvm.ConstNull(e.LLVMType())}, false))
}

func (g *moduleGenerator) getTypeSize(t llvm.Type) int {
	return int(llvm.NewTargetData(g.module.DataLayout()).TypeAllocSize(t))
}
