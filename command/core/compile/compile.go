package compile

import (
	"github.com/raviqqe/lazy-ein/command/core/ast"
	"github.com/raviqqe/lazy-ein/command/core/compile/canonicalize"
	"github.com/raviqqe/lazy-ein/command/core/validate"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
)

// Compile compiles a module into LLVM IR.
func Compile(m ast.Module) (llvm.Module, error) {
	if err := validate.Validate(m); err != nil {
		return llvm.Module{}, err
	}

	return newModuleGenerator().Generate(canonicalize.Canonicalize(m))
}
