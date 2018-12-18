package desugar

import (
	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/compile/desugar/names"
	"github.com/ein-lang/ein/command/types"
)

func desugarApplications(m ast.Module) ast.Module {
	g := names.NewNameGenerator(m.Name() + ".application")

	return m.ConvertExpression(func(e ast.Expression) ast.Expression {
		app, ok := e.(ast.Application)

		if !ok {
			return e
		}

		bs := make([]ast.Bind, 0, len(app.Arguments())+1)
		f, ok := app.Function().(ast.Variable)

		if !ok {
			f = ast.NewVariable(g.Generate("function"))
			bs = append(bs, ast.NewBind(f.Name(), types.NewVariable(nil), app.Function()))
		}

		as := make([]ast.Expression, 0, len(app.Arguments()))

		for _, a := range app.Arguments() {
			v, ok := a.(ast.Variable)

			if !ok {
				v = ast.NewVariable(g.Generate("argument"))
				bs = append(bs, ast.NewBind(v.Name(), types.NewVariable(nil), a))
			}

			as = append(as, v)
		}

		if len(bs) == 0 {
			return app
		}

		return ast.NewLet(bs, ast.NewApplication(f, as))
	}).(ast.Module)
}
