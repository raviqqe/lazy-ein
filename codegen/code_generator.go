package codegen

import "llvm.org/llvm/bindings/go/llvm"

// CodeGenerator is a code generator.
type CodeGenerator struct {
	module llvm.Module
}

// NewCodeGenerator creates a new code generator.
func NewCodeGenerator(s string) *CodeGenerator {
	return &CodeGenerator{llvm.NewModule(s)}
}

// Generate generates codes.
func (g *CodeGenerator) generate() llvm.Module {
	return g.module
}
