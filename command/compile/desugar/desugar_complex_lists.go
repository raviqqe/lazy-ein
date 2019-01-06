package desugar

import (
	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/compile/desugar/names"
	"github.com/ein-lang/ein/command/types"
)

func desugarComplexLists(m ast.Module) ast.Module {
	g := names.NewNameGenerator("list")

	return m.ConvertExpressions(func(e ast.Expression) ast.Expression {
		l, ok := e.(ast.List)

		if !ok {
			return e
		}

		bs := make([]ast.Bind, 0, len(l.Elements()))
		es := make([]ast.Expression, 0, len(l.Elements()))

		for _, e := range l.Elements() {
			v, ok := e.(ast.Variable)

			if !ok {
				v = ast.NewVariable(g.Generate("element"))
				bs = append(bs, ast.NewBind(v.Name(), types.NewUnknown(nil), e))
			}

			es = append(es, v)
		}

		if len(bs) == 0 {
			return l
		}

		return ast.NewLet(bs, ast.NewList(l.Type(), es))
	}).(ast.Module)
}
