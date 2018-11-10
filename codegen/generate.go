package codegen

import (
	"github.com/raviqqe/stg/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

// Generate generates a code for a module.
func Generate(m ast.Module) (llvm.Module, error) {
	l := llvm.NewModule(m.Name())

	if err := newModuleGenerator(l).Generate(m.Binds()); err != nil {
		return llvm.Module{}, err
	}

	return l, nil
}
