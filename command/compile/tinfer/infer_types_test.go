package tinfer_test

import (
	"testing"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/compile/tinfer"
	"github.com/ein-lang/ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestInferTypesWithLetExpressions(t *testing.T) {
	for _, ls := range [][2]ast.Let{
		// Constant expressions
		{
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"x",
						types.NewVariable(nil),
						ast.NewNumber(42),
					),
				},
				ast.NewNumber(42),
			),
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"x",
						inferredVariable(t, types.NewNumber(nil)),
						ast.NewNumber(42),
					),
				},
				ast.NewNumber(42),
			),
		},
		// Nested let expressions
		{
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"x",
						types.NewVariable(nil),
						ast.NewLet(
							[]ast.Bind{
								ast.NewBind(
									"x",
									types.NewVariable(nil),
									ast.NewNumber(42),
								),
							},
							ast.NewNumber(42),
						),
					),
				},
				ast.NewNumber(42),
			),
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"x",
						inferredVariable(t, types.NewNumber(nil)),
						ast.NewLet(
							[]ast.Bind{
								ast.NewBind(
									"x",
									inferredVariable(t, types.NewNumber(nil)),
									ast.NewNumber(42),
								),
							},
							ast.NewNumber(42),
						),
					),
				},
				ast.NewNumber(42),
			),
		},
		// Local variables
		{
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"x",
						types.NewVariable(nil),
						ast.NewLet(
							[]ast.Bind{
								ast.NewBind(
									"y",
									types.NewVariable(nil),
									ast.NewVariable("x"),
								),
							},
							ast.NewNumber(42),
						),
					),
				},
				ast.NewNumber(42),
			),
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"x",
						inferredVariable(t, types.NewNumber(nil)),
						ast.NewLet(
							[]ast.Bind{
								ast.NewBind(
									"y",
									inferredVariable(t, inferredVariable(t, types.NewNumber(nil))),
									ast.NewVariable("x"),
								),
							},
							ast.NewNumber(42),
						),
					),
				},
				ast.NewNumber(42),
			),
		},
		// Mutually recursive binds
		{
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind("x", types.NewVariable(nil), ast.NewVariable("y")),
					ast.NewBind("y", types.NewVariable(nil), ast.NewNumber(42)),
				},
				ast.NewNumber(42),
			),
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"x",
						inferredVariable(t, inferredVariable(t, types.NewNumber(nil))),
						ast.NewVariable("y"),
					),
					ast.NewBind("y", inferredVariable(t, types.NewNumber(nil)), ast.NewNumber(42)),
				},
				ast.NewNumber(42),
			),
		},
		// Functions with single arguments
		{
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewVariable(nil),
						ast.NewLambda([]string{"x"}, ast.NewNumber(42)),
					),
				},
				ast.NewNumber(42),
			),
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"f",
						inferredVariable(
							t,
							types.NewFunction(
								inferredVariable(t, nil), // TODO: Infer a real type.
								types.NewNumber(nil),
								nil,
							),
						),
						ast.NewLambda([]string{"x"}, ast.NewNumber(42)),
					),
				},
				ast.NewNumber(42),
			),
		},
		// Functions with multiple arguments
		{
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewVariable(nil),
						ast.NewLambda([]string{"x", "y"}, ast.NewNumber(42)),
					),
				},
				ast.NewNumber(42),
			),
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"f",
						inferredVariable(
							t,
							types.NewFunction(
								inferredVariable(t, nil), // TODO: Infer a real type.
								types.NewFunction(
									inferredVariable(t, nil), // TODO: Infer a real type.
									types.NewNumber(nil),
									nil,
								),
								nil,
							),
						),
						ast.NewLambda([]string{"x", "y"}, ast.NewNumber(42)),
					),
				},
				ast.NewNumber(42),
			),
		},
	} {
		m, err := tinfer.InferTypes(
			ast.NewModule("", []ast.Bind{ast.NewBind("bar", types.NewNumber(nil), ls[0])}),
		)

		assert.Nil(t, err)
		assert.Equal(t, ls[1], m.Binds()[0].Expression())
	}
}

func TestInferTypesWithArguments(t *testing.T) {
	m := ast.NewModule(
		"",
		[]ast.Bind{
			ast.NewBind(
				"f",
				types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
				ast.NewLambda([]string{"x"}, ast.NewVariable("x")),
			),
		},
	)
	mm, err := tinfer.InferTypes(m)

	assert.Nil(t, err)
	assert.Equal(t, m, mm)
}

func TestInferTypesWithFunctionApplications(t *testing.T) {
	m := ast.NewModule(
		"",
		[]ast.Bind{
			ast.NewBind(
				"a",
				types.NewNumber(nil),
				ast.NewNumber(42),
			),
			ast.NewBind(
				"f",
				types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
				ast.NewLambda([]string{"x"}, ast.NewVariable("x")),
			),
			ast.NewBind(
				"x",
				types.NewNumber(nil),
				ast.NewLet(
					[]ast.Bind{
						ast.NewBind(
							"y",
							types.NewVariable(nil),
							ast.NewApplication(ast.NewVariable("f"), []ast.Expression{ast.NewVariable("a")}),
						),
					},
					ast.NewVariable("y"),
				),
			),
		},
	)
	mm, err := tinfer.InferTypes(m)

	assert.Nil(t, err)
	assert.Equal(
		t,
		inferredVariable(t, types.NewNumber(nil)),
		mm.Binds()[2].Expression().(ast.Let).Binds()[0].Type(),
	)
}

func TestInferTypesErrorWithUnknownVarabiles(t *testing.T) {
	_, err := tinfer.InferTypes(
		ast.NewModule(
			"",
			[]ast.Bind{
				ast.NewBind(
					"x",
					types.NewNumber(nil),
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind("y", types.NewVariable(nil), ast.NewVariable("z")),
						},
						ast.NewNumber(42),
					),
				),
			},
		),
	)

	assert.Error(t, err)
}

func inferredVariable(t *testing.T, tt types.Type) *types.Variable {
	v := types.NewVariable(nil)

	assert.Nil(t, v.Unify(tt))

	return v
}
