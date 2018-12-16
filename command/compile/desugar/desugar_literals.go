package desugar

import (
	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/compile/desugar/names"
	"github.com/ein-lang/ein/command/types"
)

func desugarLiterals(m ast.Module) ast.Module {
	g := names.NewNameGenerator(m.Name())
	bs := []ast.Bind{}

	for _, b := range m.Binds() {
		if l, ok := b.Expression().(ast.Literal); ok && len(b.Arguments()) == 0 {
			bs = append(
				bs,
				ast.NewBind(
					b.Name(),
					nil,
					types.NewUnboxed(b.Type(), b.Type().DebugInformation()),
					ast.NewUnboxed(l),
				),
			)

			continue
		}

		bs = append(bs, b.ConvertExpression(func(e ast.Expression) ast.Expression {
			l, ok := e.(ast.Literal)

			if !ok {
				return e
			}

			s := g.Generate("literal")

			// TODO: Handle other literals.
			switch l := l.(type) {
			case ast.Number:
				bs = append(
					bs,
					ast.NewBind(s, nil, types.NewUnboxed(types.NewNumber(nil), nil), ast.NewUnboxed(l)),
				)
				return ast.NewVariable(s)
			}

			panic("unreachable")
		}).(ast.Bind))
	}

	return ast.NewModule(m.Name(), bs)
}
