package validate

import (
	"testing"

	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateError(t *testing.T) {
	tt := types.NewAlgebraic(types.NewConstructor())

	assert.Error(
		t,
		checkRecursiveBinds(
			ast.NewModule(
				nil,
				[]ast.Bind{
					ast.NewBind(
						"x",
						ast.NewVariableLambda(
							nil,
							true,
							ast.NewFunctionApplication(ast.NewVariable("x"), nil),
							types.NewBoxed(tt),
						),
					),
				},
			),
		),
	)
}

func TestValidateErrorWithLetExpressions(t *testing.T) {
	tt := types.NewAlgebraic(types.NewConstructor())

	assert.Error(
		t,
		checkRecursiveBinds(
			ast.NewModule(
				nil,
				[]ast.Bind{
					ast.NewBind(
						"x",
						ast.NewVariableLambda(
							nil,
							true,
							ast.NewLet(
								[]ast.Bind{
									ast.NewBind(
										"y",
										ast.NewVariableLambda(
											[]ast.Argument{ast.NewArgument("y", types.NewBoxed(tt))},
											true,
											ast.NewFunctionApplication(ast.NewVariable("y"), nil),
											types.NewBoxed(tt),
										),
									),
								},
								ast.NewFunctionApplication(ast.NewVariable("y"), nil),
							),
							types.NewBoxed(tt),
						),
					),
				},
			),
		),
	)
}
