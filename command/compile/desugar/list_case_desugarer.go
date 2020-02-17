package desugar

import (
	"github.com/raviqqe/lazy-ein/command/ast"
	"github.com/raviqqe/lazy-ein/command/compile/desugar/names"
	"github.com/raviqqe/lazy-ein/command/types"
)

type listCaseDesugarer struct {
	nameGenerator names.NameGenerator
}

func newListCaseDesugarer() listCaseDesugarer {
	return listCaseDesugarer{names.NewNameGenerator("list-case")}
}

func (dd listCaseDesugarer) Desugar(c ast.Case) ast.Expression {
	c = desugarHiddenDefaultAlternative(c)

	if len(findNonEmptyListAlternatives(c.Alternatives())) == 0 {
		return c
	}

	d, ok := c.DefaultAlternative()

	if !ok {
		return ast.NewCaseWithoutDefault(c.Argument(), c.Type(), dd.createListAlternatives(c, nil))
	} else if d.Variable() == "" {
		return ast.NewCase(
			c.Argument(),
			types.NewUnknown(nil),
			dd.createListAlternatives(c, d.Expression()),
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
		[]ast.Bind{ast.NewBind(s, types.NewUnknown(nil), c.Argument())},
		ast.NewCase(
			ast.NewVariable(s),
			types.NewUnknown(nil),
			dd.createListAlternatives(c, e),
			ast.NewDefaultAlternative("", e),
		),
	)
}

func (dd listCaseDesugarer) createListAlternatives(
	c ast.Case,
	defaultExpr ast.Expression,
) []ast.Alternative {
	as := make([]ast.Alternative, 0, 2)

	if a, ok := dd.desugarNonEmptyListAlternatives(c.Alternatives(), defaultExpr); ok {
		as = append(as, a)
	}

	if a, ok := findEmptyListAlternative(c.Alternatives()); ok {
		as = append(as, a)
	}

	return as
}

func (dd listCaseDesugarer) desugarNonEmptyListAlternatives(
	as []ast.Alternative,
	defaultExpr ast.Expression,
) (ast.Alternative, bool) {
	as = findNonEmptyListAlternatives(as)

	if len(as) == 0 {
		return ast.Alternative{}, false
	}

	head := dd.nameGenerator.Generate("head")
	tail := dd.nameGenerator.Generate("tail")

	return ast.NewAlternative(
		ast.NewList(
			as[0].Pattern().(ast.List).Type(),
			[]ast.ListArgument{
				ast.NewListArgument(ast.NewVariable(head), false),
				ast.NewListArgument(ast.NewVariable(tail), true),
			},
		),
		dd.createHeadCase(head, tail, as, defaultExpr),
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

func (dd listCaseDesugarer) createHeadCase(head, tail string, as []ast.Alternative, defaultExpr ast.Expression) ast.Expression {
	if len(as) == 0 {
		return defaultExpr
	}

	ass := [][]ast.Alternative{}

	for i, a := range as {
		if v, ok := headPattern(a).(ast.Variable); ok {
			defaultExpr = dd.Desugar(
				ast.NewCase(
					ast.NewVariable(tail),
					a.Pattern().(ast.List).Type(),
					[]ast.Alternative{
						ast.NewAlternative(
							tailPattern(a),
							ast.NewLet(
								[]ast.Bind{ast.NewBind(v.Name(), types.NewUnknown(nil), ast.NewVariable(head))},
								a.Expression(),
							),
						),
					},
					ast.NewDefaultAlternative(
						"",
						dd.createHeadCase(head, tail, as[i+1:], defaultExpr),
					),
				),
			)

			return ast.NewCase(
				ast.NewVariable(head),
				types.NewUnknown(nil),
				dd.createHeadAlternatives(tail, ass, defaultExpr),
				ast.NewDefaultAlternative("", defaultExpr),
			)
		} else if len(ass) > 0 && ast.PatternsEqual(headPattern(a), headPattern(ass[len(ass)-1][0])) {
			ass[len(ass)-1] = append(ass[len(ass)-1], a)
			continue
		}

		ass = append(ass, []ast.Alternative{a})
	}

	return ast.NewCase(
		ast.NewVariable(head),
		types.NewUnknown(nil),
		dd.createHeadAlternatives(tail, ass, defaultExpr),
		ast.NewDefaultAlternative("", defaultExpr),
	)
}

func (dd listCaseDesugarer) createHeadAlternatives(
	tail string,
	ass [][]ast.Alternative,
	defaultExpr ast.Expression,
) []ast.Alternative {
	as := make([]ast.Alternative, 0, len(ass))

	for _, aas := range ass {
		as = append(
			as,
			ast.NewAlternative(
				headPattern(aas[0]),
				dd.Desugar(
					ast.NewCase(
						ast.NewVariable(tail),
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

	return ast.NewCase(c.Argument(), c.Type(), as, d)
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
