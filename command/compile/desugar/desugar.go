package desugar

import "github.com/raviqqe/jsonxx/command/ast"

// Desugar desugars an AST.
func Desugar(m ast.Module) ast.Module {
	return desugarLiterals(m)
}
