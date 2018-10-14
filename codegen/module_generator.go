package codegen

import (
	"github.com/raviqqe/stg/ast"
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
		f := llvm.AddFunction(
			g.module,
			toEntryName(b.Name()),
			llvm.FunctionType(
				llvm.DoubleType(),
				append([]llvm.Type{voidPointerType}),
				false,
			),
		)
		newFunctionGenerator(f).Generate(b.Lambda().Body())

		if err := llvm.VerifyFunction(f, llvm.AbortProcessAction); err != nil {
			return err
		}

		g.createClosure(b.Name(), f)
	}

	return llvm.VerifyModule(g.module, llvm.AbortProcessAction)
}

func (g *moduleGenerator) createClosure(s string, f llvm.Value) {
	llvm.AddGlobal(
		g.module,
		llvm.StructType([]llvm.Type{f.Type()}, false),
		s,
	).SetInitializer(llvm.ConstStruct([]llvm.Value{f}, false))
}
