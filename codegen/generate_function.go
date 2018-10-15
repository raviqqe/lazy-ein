package codegen

import (
	"github.com/raviqqe/stg/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

func generateFunction(b ast.Bind, m llvm.Module) llvm.Value {
	g := newFunctionGenerator(b, m)
	g.Generate()
	return g.function
}
