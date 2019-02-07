package desugar

import (
	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/types"
)

const (
	headVariable = "$head"
	tailVariable = "$tail"
)

// desugarListCases converts all elements of list patterns into variables.
func desugarListCases(m ast.Module) ast.Module {
	return m.ConvertExpressions(mayDesugarListCase).(ast.Module)
}

func mayDesugarListCase(e ast.Expression) ast.Expression {
	c, ok := e.(ast.Case)

	if !ok || len(c.Alternatives()) == 0 {
		return e
	} else if _, ok := c.Alternatives()[0].Pattern().(ast.List); !ok {
		return e
	}

	return desugarListCase(c)
}

func desugarListCase(c ast.Case) ast.Expression {
	c = desugarHiddenDefaultAlternative(c)

	if len(findNonEmptyListAlternatives(c.Alternatives())) == 0 {
		return c
	}

	d, ok := c.DefaultAlternative()

	if !ok {
		return ast.NewCaseWithoutDefault(c.Expression(), c.Type(), createListAlternatives(c, nil))
	} else if d.Variable() == "" {
		return ast.NewCase(
			c.Expression(),
			types.NewUnknown(nil),
			createListAlternatives(c, d.Expression()),
			d,
		)
	}

	// TODO: Prove that name generation is not necessary here.
	s := "$default-alternative." + d.Variable()
	e := ast.NewLet(
		[]ast.Bind{ast.NewBind(d.Variable(), c.Type(), ast.NewVariable(s))},
		d.Expression(),
	)

	return ast.NewLet(
		[]ast.Bind{ast.NewBind(s, types.NewUnknown(nil), c.Expression())},
		ast.NewCase(
			ast.NewVariable(s),
			types.NewUnknown(nil),
			createListAlternatives(c, e),
			ast.NewDefaultAlternative("", e),
		),
	)
}

func createListAlternatives(c ast.Case, defaultExpr ast.Expression) []ast.Alternative {
	as := make([]ast.Alternative, 0, 2)

	if a, ok := desugarNonEmptyListAlternatives(c.Alternatives(), defaultExpr); ok {
		as = append(as, a)
	}

	if a, ok := findEmptyListAlternative(c.Alternatives()); ok {
		as = append(as, a)
	}

	return as
}

func desugarNonEmptyListAlternatives(
	as []ast.Alternative,
	defaultExpr ast.Expression,
) (ast.Alternative, bool) {
	as = findNonEmptyListAlternatives(as)

	if len(as) == 0 {
		return ast.Alternative{}, false
	}

	return ast.NewAlternative(
		ast.NewList(
			as[0].Pattern().(ast.List).Type(),
			[]ast.ListArgument{
				ast.NewListArgument(ast.NewVariable(headVariable), false),
				ast.NewListArgument(ast.NewVariable(tailVariable), true),
			},
		),
		createHeadCase(as, defaultExpr),
	), true
}

func findNonEmptyListAlternatives(as []ast.Alternative) []ast.Alternative {
	aas := make([]ast.Alternative, 0, len(as))

	for _, a := range as {
		if len(a.Pattern().(ast.List).Arguments()) != 0 {
			aas = append(aas, a)
		}
	}

	return aas
}

func findEmptyListAlternative(as []ast.Alternative) (ast.Alternative, bool) {
	for _, a := range as {
		if len(a.Pattern().(ast.List).Arguments()) == 0 {
			return a, true
		}
	}

	return ast.Alternative{}, false
}

// TODO
func createHeadCase(as []ast.Alternative, defaultExpr ast.Expression) ast.Expression {
	ass := [][]ast.Alternative{}

	for i, a := range as {
		if v, ok := headPattern(a).(ast.Variable); ok {
			e := defaultExpr

			if as := as[i+1:]; len(as) != 0 {
				e = createHeadCase(as, defaultExpr)
			}

			return ast.NewCase(
				ast.NewVariable(headVariable),
				types.NewUnknown(nil),
				createHeadAlternatives(ass, defaultExpr),
				ast.NewDefaultAlternative(
					v.Name(),
					ast.NewCase(
						ast.NewVariable(tailVariable),
						a.Pattern().(ast.List).Type(),
						[]ast.Alternative{ast.NewAlternative(tailPattern(a), a.Expression())},
						ast.NewDefaultAlternative("", e),
					),
				),
			)
		} else if len(ass) > 0 && ast.PatternsEqual(headPattern(a), headPattern(ass[len(ass)-1][0])) {
			ass[len(ass)-1] = append(ass[len(ass)-1], a)
			continue
		}

		ass = append(ass, []ast.Alternative{a})
	}

	return ast.NewCase(
		ast.NewVariable(headVariable),
		types.NewUnknown(nil),
		createHeadAlternatives(ass, defaultExpr),
		ast.NewDefaultAlternative("", defaultExpr),
	)
}

func createHeadAlternatives(
	ass [][]ast.Alternative,
	defaultExpr ast.Expression,
) []ast.Alternative {
	as := make([]ast.Alternative, 0, len(ass))

	for _, aas := range ass {
		as = append(
			as,
			ast.NewAlternative(
				headPattern(aas[0]),
				desugarListCase(
					ast.NewCase(
						ast.NewVariable(tailVariable),
						types.NewUnknown(nil),
						createTailAlternatives(aas),
						ast.NewDefaultAlternative("", defaultExpr),
					),
				),
			),
		)
	}

	return as
}

func createTailAlternatives(as []ast.Alternative) []ast.Alternative {
	aas := make([]ast.Alternative, 0, len(as))

	for _, a := range as {
		aas = append(aas, ast.NewAlternative(tailPattern(a), a.Expression()))
	}

	return aas
}

func desugarHiddenDefaultAlternative(c ast.Case) ast.Case {
	d, as, ok := findHiddenDefaultAlternative(c.Alternatives())

	if !ok {
		return c
	}

	return ast.NewCase(c.Expression(), c.Type(), as, d)
}

func findHiddenDefaultAlternative(
	as []ast.Alternative,
) (ast.DefaultAlternative, []ast.Alternative, bool) {
	for i, a := range as {
		if l, ok := a.Pattern().(ast.List); ok {
			if len(l.Arguments()) == 1 && l.Arguments()[0].Expanded() {
				if v, ok := l.Arguments()[0].Expression().(ast.Variable); ok {
					return ast.NewDefaultAlternative(v.Name(), a.Expression()), append(as[:i], as[i+1:]...), true
				}
			}
		}
	}

	return ast.DefaultAlternative{}, nil, false
}

func headPattern(a ast.Alternative) ast.Expression {
	return a.Pattern().(ast.List).Arguments()[0].Expression()
}

func tailPattern(a ast.Alternative) ast.Expression {
	l := a.Pattern().(ast.List)
	args := l.Arguments()[1:]

	if len(args) == 0 {
		return ast.NewList(l.Type(), nil)
	}

	return ast.NewList(l.Type(), args)
}
