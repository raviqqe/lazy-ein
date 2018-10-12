package codegen

import (
	"github.com/raviqqe/stg/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

// CodeGenerator is a code generator.
type CodeGenerator struct {
	module llvm.Module
}

// NewCodeGenerator creates a new code generator.
func NewCodeGenerator(s string) *CodeGenerator {
	return &CodeGenerator{llvm.NewModule(s)}
}

// GenerateModule generates codes for modules.
func (g *CodeGenerator) GenerateModule(bs []ast.Bind) {
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
