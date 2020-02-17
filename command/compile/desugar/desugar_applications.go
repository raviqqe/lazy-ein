package desugar

import (
	"github.com/raviqqe/lazy-ein/command/ast"
	"github.com/raviqqe/lazy-ein/command/compile/desugar/names"
	"github.com/raviqqe/lazy-ein/command/types"
)

func desugarApplications(m ast.Module) ast.Module {
	g := names.NewNameGenerator("application")

	return m.ConvertExpressions(func(e ast.Expression) ast.Expression {
		app, ok := e.(ast.Application)

		if !ok {
			return e
		}

		bs := make([]ast.Bind, 0, len(app.Arguments())+1)
		f, ok := app.Function().(ast.Variable)

		if !ok {
			f = ast.NewVariable(g.Generate("function"))
			bs = append(bs, ast.NewBind(f.Name(), types.NewUnknown(nil), app.Function()))
		}

		as := make([]ast.Expression, 0, len(app.Arguments()))

		for _, a := range app.Arguments() {
			v, ok := a.(ast.Variable)

			if !ok {
				v = ast.NewVariable(g.Generate("argument"))
				bs = append(bs, ast.NewBind(v.Name(), types.NewUnknown(nil), a))
			}

			as = append(as, v)
		}

		if len(bs) == 0 {
			return app
		}

		return ast.NewLet(bs, ast.NewApplication(f, as))
	})
}
