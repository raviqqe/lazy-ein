package codegen

import (
	"github.com/raviqqe/stg/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

// Generate generates a code for a module.
func Generate(s string, bs []ast.Bind) (llvm.Module, error) {
	g := newModuleGenerator(s)

	if err := g.Generate(bs); err != nil {
		return llvm.Module{}, err
	}

	return g.module, nil
}
