package codegen

import (
	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/types"
	"llvm.org/llvm/bindings/go/llvm"
)

func generateFunction(b ast.Bind, e types.Environment, m llvm.Module) llvm.Value {
	g := newFunctionGenerator(b, e, m)
	g.Generate()
	return g.function
}
