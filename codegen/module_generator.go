package codegen

import (
	"github.com/raviqqe/stg/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

// ModuleGenerator is a code generator.
type ModuleGenerator struct {
	module llvm.Module
}

// NewModuleGenerator creates a new code generator.
func NewModuleGenerator(s string) *ModuleGenerator {
	return &ModuleGenerator{llvm.NewModule(s)}
}

// Generate generates a module.
func (g *ModuleGenerator) Generate(bs []ast.Bind) {
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

		g.createClosure(b.Name(), f)
	}
}

func (g *ModuleGenerator) createClosure(s string, f llvm.Value) {
	llvm.AddGlobal(
		g.module,
		llvm.StructType([]llvm.Type{f.Type()}, false),
		s,
	).SetInitializer(llvm.ConstStruct([]llvm.Value{f}, false))
}
