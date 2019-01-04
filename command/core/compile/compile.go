package compile

import (
	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/validate"
	"llvm.org/llvm/bindings/go/llvm"
)

// Compile compiles a module into LLVM IR.
func Compile(m ast.Module) (llvm.Module, error) {
	if err := validate.Validate(m); err != nil {
		return llvm.Module{}, err
	}

	mm := llvm.NewModule(m.Name())
	g, err := newModuleGenerator(mm, m.TypeDefinitions())

	if err != nil {
		return llvm.Module{}, err
	} else if err := g.Generate(m.Binds()); err != nil {
		return llvm.Module{}, err
	}

	return mm, nil
}
