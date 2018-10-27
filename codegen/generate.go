package codegen

import (
	"github.com/raviqqe/stg/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

// Generate generates a code for a module.
func Generate(s string, bs []ast.Bind) (llvm.Module, error) {
	m := llvm.NewModule(s)

	if err := newModuleGenerator(m).Generate(bs); err != nil {
		return llvm.Module{}, err
	}

	return m, nil
}
