package codegen

import (
	"github.com/raviqqe/stg/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

// Generate generates a code for a module.
func Generate(m ast.Module) (llvm.Module, error) {
	mm := llvm.NewModule(m.Name())
	g, err := newModuleGenerator(mm, m.ConstructorDefinitions())

	if err != nil {
		return llvm.Module{}, err
	} else if err := g.Generate(m.Binds()); err != nil {
		return llvm.Module{}, err
	}

	return mm, nil
}
