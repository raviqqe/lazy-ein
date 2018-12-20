package desugar

import "github.com/ein-lang/ein/command/ast"

// Desugar desugars an AST.
func Desugar(m ast.Module) ast.Module {
	for _, f := range []func(ast.Module) ast.Module{
		desugarLiterals,
		desugarComplexApplications,
	} {
		m = f(m)
	}

	return m
}
