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

	return corecompile.Compile(renameGlobalVariables(m, mm))
}

func compileToCore(m ast.Module) (coreast.Module, error) {
	m, err := tinfer.InferTypes(desugar.WithoutTypes(m), nil)

	if err != nil {
		return coreast.Module{}, err
	}

	return newCompiler().Compile(desugar.WithTypes(m))
}

func renameGlobalVariables(m ast.Module, mm coreast.Module) coreast.Module {
	vs := make(map[string]string, len(mm.Binds()))
	bs := make([]coreast.Bind, 0, len(mm.Binds()))

	for _, b := range mm.Binds() {
		s := string(m.Name()) + "." + b.Name()

		vs[b.Name()] = s

		if b.Name() == ast.MainFunctionName {
			bs = append(bs, coreast.NewBind("ein_main", b.Lambda()))
		} else {
			bs = append(bs, coreast.NewBind(s, b.Lambda()))
		}
	}

	return coreast.NewModule(mm.Declarations(), bs).RenameVariables(vs)
}
