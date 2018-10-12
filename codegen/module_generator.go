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

// GenerateModule generates codes for modules.
func (g *ModuleGenerator) GenerateModule(bs []ast.Bind) {
	for _, b := range bs {
		t := llvm.FunctionType(
			llvm.DoubleType(),
			append(
				[]llvm.Type{voidPointerType},
				thunkPointerArrayType(len(b.Lambda().Arguments()))...,
			),
			false,
		)

		f := llvm.AddFunction(g.module, toEntryName(b.Name()), t)
		newFunctionGenerator(f).Generate(b.Lambda().Body())

		llvm.AddGlobal(
			g.module,
			llvm.StructType(
				[]llvm.Type{
					t,
					llvm.ArrayType(thunkPointerType, len(b.Lambda().FreeVariables())),
				},
				false,
			),
			b.Name(),
		)
	}
}
