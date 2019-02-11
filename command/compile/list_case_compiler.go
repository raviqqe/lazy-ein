package compile

import (
	"github.com/ein-lang/ein/command/ast"
	coreast "github.com/ein-lang/ein/command/core/ast"
	coretypes "github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/types"
)

type listCaseCompiler struct {
	compiler
}

func newListCaseCompiler(c compiler) listCaseCompiler {
	return listCaseCompiler{c}
}

func (c listCaseCompiler) Compile(cc ast.Case) (coreast.Expression, error) {
	t := coretypes.Unbox(cc.Type().ToCore()).(coretypes.Algebraic)
	arg, err := c.compileExpression(cc.Argument())

	if err != nil {
		return nil, err
	}

	as := make([]coreast.AlgebraicAlternative, 0, 2)

	for _, a := range cc.Alternatives() {
		l := a.Pattern().(ast.List)

		if len(l.Arguments()) != 0 {
			c = c.addVariable(
				l.Arguments()[0].Expression().(ast.Variable).Name(),
				cc.Type().(types.List).Element().ToCore(),
			).addVariable(
				l.Arguments()[1].Expression().(ast.Variable).Name(),
				cc.Type().ToCore(),
			)
		}

		e, err := c.compileExpression(a.Expression())

		if err != nil {
			return nil, err
		}

		aa := coreast.NewAlgebraicAlternative(coreast.NewConstructor(t, 1), nil, e)

		if len(l.Arguments()) != 0 {
			aa = coreast.NewAlgebraicAlternative(
				coreast.NewConstructor(t, 0),
				[]string{
					l.Arguments()[0].Expression().(ast.Variable).Name(),
					l.Arguments()[1].Expression().(ast.Variable).Name(),
				},
				e,
			)
		}

		as = append(as, aa)
	}

	d, ok := cc.DefaultAlternative()

	if !ok {
		return coreast.NewAlgebraicCaseWithoutDefault(arg, as), nil
	}

	e, err := c.addVariable(d.Variable(), coretypes.NewBoxed(t)).compileExpression(
		d.Expression(),
	)

	if err != nil {
		return nil, err
	}

	vs, err := c.compileFreeVariables(cc.Argument())

	if err != nil {
		return nil, err
	}

	// TODO: Prove that name generation is not necessary here.
	s := "$argument." + d.Variable()

	return coreast.NewLet(
		[]coreast.Bind{
			coreast.NewBind(s, coreast.NewVariableLambda(vs, true, arg, coretypes.NewBoxed(t))),
		},
		coreast.NewAlgebraicCase(
			coreast.NewFunctionApplication(coreast.NewVariable(s), nil),
			as,
			coreast.NewDefaultAlternative(
				"",
				coreast.NewLet(
					[]coreast.Bind{
						coreast.NewBind(
							d.Variable(),
							coreast.NewVariableLambda(
								[]coreast.Argument{coreast.NewArgument(s, coretypes.NewBoxed(t))},
								true,
								coreast.NewFunctionApplication(coreast.NewVariable(s), nil),
								coretypes.NewBoxed(t),
							),
						),
					},
					e,
				),
			),
		),
	), nil
}

func (c listCaseCompiler) addVariable(s string, t coretypes.Type) listCaseCompiler {
	return listCaseCompiler{c.compiler.addVariable(s, t)}
}
