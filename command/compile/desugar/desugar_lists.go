package desugar

import (
	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/compile/desugar/names"
	"github.com/ein-lang/ein/command/types"
)

func desugarLists(m ast.Module) ast.Module {
	g := names.NewNameGenerator("list")

	return m.ConvertExpressions(func(e ast.Expression) ast.Expression {
		l, ok := e.(ast.List)

		if !ok {
			return e
		}

		bs := make([]ast.Bind, 0, len(l.Arguments()))
		as := make([]ast.ListArgument, 0, len(l.Arguments()))

		for _, a := range l.Arguments() {
			if a.Expanded() {
				panic("not implemented")
			} else if _, ok := a.Expression().(ast.Variable); ok {
				as = append(as, a)
				continue
			}

			s := g.Generate("element")
			bs = append(bs, ast.NewBind(s, types.NewUnknown(nil), a.Expression()))
			as = append(as, ast.NewListArgument(ast.NewVariable(s), false))
		}

		if len(bs) == 0 {
			return l
		}

		return ast.NewLet(bs, ast.NewList(l.Type(), as))
	}).(ast.Module)
}
