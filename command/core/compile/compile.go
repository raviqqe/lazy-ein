package compile

import (
	"github.com/raviqqe/jsonxx/command/core/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

// Compile compiles a module into LLVM IR.
func Compile(m ast.Module) (llvm.Module, error) {
	mm := llvm.NewModule(m.Name())
	g, err := newModuleGenerator(mm, m.ConstructorDefinitions())

	if err != nil {
		return llvm.Module{}, err
	} else if err := g.Generate(m.Binds()); err != nil {
		return llvm.Module{}, err
	}

	return mm, nil
}
