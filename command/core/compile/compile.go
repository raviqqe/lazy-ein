package compile

import (
	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/compile/canonicalize"
	"github.com/ein-lang/ein/command/core/validate"
	"llvm.org/llvm/bindings/go/llvm"
)

// Compile compiles a module into LLVM IR.
func Compile(m ast.Module) (llvm.Module, error) {
	if err := validate.Validate(m); err != nil {
		return llvm.Module{}, err
	}

	m = canonicalize.Canonicalize(m)

	mm := llvm.NewModule(m.Name())
	g, err := newModuleGenerator(mm, m)

	if err != nil {
		return llvm.Module{}, err
	} else if err := g.Generate(m.Binds()); err != nil {
		return llvm.Module{}, err
	}

	return mm, nil
}
