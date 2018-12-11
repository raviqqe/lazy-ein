package desugar

import "github.com/ein-lang/ein/command/ast"

// Desugar desugars an AST.
func Desugar(m ast.Module) ast.Module {
	return desugarLiterals(m)
}
