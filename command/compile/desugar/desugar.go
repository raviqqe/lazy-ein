package desugar

import "github.com/raviqqe/lazy-ein/command/ast"

// WithoutTypes desugars an AST without type information.
func WithoutTypes(m ast.Module) ast.Module {
	for _, f := range []func(ast.Module) ast.Module{
		desugarLiterals,
		desugarApplications,
		desugarBinaryOperations,
		desugarLists,
		desugarListCases,
	} {
		m = f(m)
	}

	return m
}

// WithTypes desugars an AST with type information.
func WithTypes(m ast.Module) ast.Module {
	return desugarPartialApplications(m)
}
