package desugar

import "github.com/raviqqe/lazy-ein/command/ast"

// desugarListCases converts all elements of list patterns into variables.
func desugarListCases(m ast.Module) ast.Module {
	return m.ConvertExpressions(mayDesugarListCase)
}

func mayDesugarListCase(e ast.Expression) ast.Expression {
	c, ok := e.(ast.Case)

	if !ok || len(c.Alternatives()) == 0 {
		return e
	} else if _, ok := c.Alternatives()[0].Pattern().(ast.List); !ok {
		return e
	}

	return newListCaseDesugarer().Desugar(c)
}
