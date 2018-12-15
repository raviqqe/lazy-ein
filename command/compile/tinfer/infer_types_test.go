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
						nil,
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
						nil,
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
						nil,
						types.NewVariable(nil),
						ast.NewLet(
							[]ast.Bind{
								ast.NewBind(
									"x",
									nil,
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
						nil,
						inferredVariable(t, types.NewNumber(nil)),
						ast.NewLet(
							[]ast.Bind{
								ast.NewBind(
									"x",
									nil,
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
						nil,
						types.NewVariable(nil),
						ast.NewLet(
							[]ast.Bind{
								ast.NewBind(
									"y",
									nil,
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
						nil,
						inferredVariable(t, types.NewNumber(nil)),
						ast.NewLet(
							[]ast.Bind{
								ast.NewBind(
									"y",
									nil,
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
					ast.NewBind("x", nil, types.NewVariable(nil), ast.NewVariable("y")),
					ast.NewBind("y", nil, types.NewVariable(nil), ast.NewNumber(42)),
				},
				ast.NewNumber(42),
			),
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"x",
						nil,
						inferredVariable(t, inferredVariable(t, types.NewNumber(nil))),
						ast.NewVariable("y"),
					),
					ast.NewBind("y", nil, inferredVariable(t, types.NewNumber(nil)), ast.NewNumber(42)),
				},
				ast.NewNumber(42),
			),
		},
		// Functions with single arguments
		{
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind("f", []string{"x"}, types.NewVariable(nil), ast.NewNumber(42)),
				},
				ast.NewNumber(42),
			),
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"f",
						[]string{"x"},
						inferredVariable(
							t,
							types.NewFunction(
								inferredVariable(t, nil), // TODO: Infer a real type.
								inferredVariable(t, types.NewNumber(nil)),
								nil,
							),
						),
						ast.NewNumber(42),
					),
				},
				ast.NewNumber(42),
			),
		},
		// Functions with multiple arguments
		{
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind("f", []string{"x", "y"}, types.NewVariable(nil), ast.NewNumber(42)),
				},
				ast.NewNumber(42),
			),
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"f",
						[]string{"x", "y"},
						inferredVariable(
							t,
							types.NewFunction(
								inferredVariable(t, nil), // TODO: Infer a real type.
								types.NewFunction(
									inferredVariable(t, nil), // TODO: Infer a real type.
									inferredVariable(t, types.NewNumber(nil)),
									nil,
								),
								nil,
							),
						),
						ast.NewNumber(42),
					),
				},
				ast.NewNumber(42),
			),
		},
	} {
		m, err := tinfer.InferTypes(
			ast.NewModule("", []ast.Bind{ast.NewBind("bar", nil, types.NewNumber(nil), ls[0])}),
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
				[]string{"x"},
				types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
				ast.NewVariable("x"),
			),
		},
	)
	mm, err := tinfer.InferTypes(m)

	assert.Nil(t, err)
	assert.Equal(t, m, mm)
}

func TestInferTypesErrorWithTooManyArguments(t *testing.T) {
	_, err := tinfer.InferTypes(
		ast.NewModule(
			"",
			[]ast.Bind{
				ast.NewBind(
					"f",
					[]string{"x"},
					types.NewNumber(nil),
					ast.NewVariable("x"),
				),
			},
		),
	)

	assert.Error(t, err)
}

func TestInferTypesErrorWithUnknownVarabiles(t *testing.T) {
	_, err := tinfer.InferTypes(
		ast.NewModule(
			"",
			[]ast.Bind{
				ast.NewBind(
					"x",
					nil,
					types.NewNumber(nil),
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind("y", nil, types.NewVariable(nil), ast.NewVariable("z")),
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
