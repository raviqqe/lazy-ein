package compile

import (
	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/compile/desugar"
	"github.com/ein-lang/ein/command/compile/tinfer"
	coreast "github.com/ein-lang/ein/command/core/ast"
	corecompile "github.com/ein-lang/ein/command/core/compile"
	"llvm.org/llvm/bindings/go/llvm"
)

// Compile compiles a module into a module in the core language.
func Compile(m ast.Module) (llvm.Module, error) {
	mm, err := compileToCore(m)

	if err != nil {
		return llvm.Module{}, err
	}

	return corecompile.Compile(mm)
}

func compileToCore(m ast.Module) (coreast.Module, error) {
	var err error
	m, err = tinfer.InferTypes(m)

	if err != nil {
		return coreast.Module{}, err
	}

	m = desugar.Desugar(m)

	return newCompiler(m).Compile(m)
}
