package codegen

import (
	"github.com/raviqqe/stg/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

// Generate generates a code for a module.
func Generate(m ast.Module) (llvm.Module, error) {
	mm := llvm.NewModule(m.Name())

	if err := newModuleGenerator(mm, m.ConstructorDefinitions()).Generate(m.Binds()); err != nil {
		return llvm.Module{}, err
	}

	return mm, nil
}
