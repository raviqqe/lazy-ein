package desugar

import (
	"github.com/raviqqe/lazy-ein/command/ast"
	"github.com/raviqqe/lazy-ein/command/compile/desugar/names"
	"github.com/raviqqe/lazy-ein/command/types"
)

func desugarBinaryOperations(m ast.Module) ast.Module {
	g := names.NewNameGenerator("binary-operation")

	return m.ConvertExpressions(func(e ast.Expression) ast.Expression {
		o, ok := e.(ast.BinaryOperation)

		if !ok {
			return e
		}

		bs := make([]ast.Bind, 0, 2)
		es := make([]ast.Expression, 0, 2)

		for _, e := range []ast.Expression{o.LHS(), o.RHS()} {
			v, ok := e.(ast.Variable)

			if !ok {
				v = ast.NewVariable(g.Generate("argument"))
				bs = append(bs, ast.NewBind(v.Name(), types.NewUnknown(nil), e))
			}

			es = append(es, v)
		}

		if len(bs) == 0 {
			return o
		}

		return ast.NewLet(bs, ast.NewBinaryOperation(o.Operator(), es[0], es[1]))
	})
}
