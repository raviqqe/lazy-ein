package desugar

import (
	"testing"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestDesugarPartialApplications(t *testing.T) {
	for _, ms := range [][2]ast.Module{
		// Empty modules
		{
			ast.NewModule("", []ast.Bind{}),
			ast.NewModule("", []ast.Bind{}),
		},
		// Variables with single arguments
		{
			ast.NewModule(
				"",
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewLambda([]string{"x"}, ast.NewVariable("x")),
					),
					ast.NewBind(
						"g",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewVariable("f"),
					),
				},
			),
			ast.NewModule(
				"",
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewLambda([]string{"x"}, ast.NewVariable("x")),
					),
					ast.NewBind(
						"g",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewLambda(
							[]string{"additional.argument-0"},
							ast.NewApplication(
								ast.NewVariable("f"),
								[]ast.Expression{ast.NewVariable("additional.argument-0")},
							),
						),
					),
				},
			),
		},
		// Variables with multiple arguments
		{
			ast.NewModule(
				"",
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewFunction(
							types.NewNumber(nil),
							types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
							nil,
						),
						ast.NewLambda([]string{"x", "y"}, ast.NewVariable("x")),
					),
					ast.NewBind(
						"g",
						types.NewFunction(
							types.NewNumber(nil),
							types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
							nil,
						),
						ast.NewVariable("f"),
					),
				},
			),
			ast.NewModule(
				"",
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewFunction(
							types.NewNumber(nil),
							types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
							nil,
						),
						ast.NewLambda([]string{"x", "y"}, ast.NewVariable("x")),
					),
					ast.NewBind(
						"g",
						types.NewFunction(
							types.NewNumber(nil),
							types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
							nil,
						),
						ast.NewLambda(
							[]string{"additional.argument-0", "additional.argument-1"},
							ast.NewApplication(
								ast.NewVariable("f"),
								[]ast.Expression{
									ast.NewVariable("additional.argument-0"),
									ast.NewVariable("additional.argument-1"),
								},
							),
						),
					),
				},
			),
		},
		// Applications
		{
			ast.NewModule(
				"",
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewFunction(
							types.NewNumber(nil),
							types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
							nil,
						),
						ast.NewLambda([]string{"x", "y"}, ast.NewVariable("x")),
					),
					ast.NewBind(
						"g",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewApplication(ast.NewVariable("f"), []ast.Expression{ast.NewVariable("x")}),
					),
				},
			),
			ast.NewModule(
				"",
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewFunction(
							types.NewNumber(nil),
							types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
							nil,
						),
						ast.NewLambda([]string{"x", "y"}, ast.NewVariable("x")),
					),
					ast.NewBind(
						"g",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewLambda(
							[]string{"additional.argument-0"},
							ast.NewApplication(
								ast.NewVariable("f"),
								[]ast.Expression{ast.NewVariable("x"), ast.NewVariable("additional.argument-0")},
							),
						),
					),
				},
			),
		},
		// Let expressions
		{
			ast.NewModule(
				"",
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewLambda([]string{"x"}, ast.NewVariable("x")),
					),
					ast.NewBind(
						"g",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewLet(
							[]ast.Bind{ast.NewBind("a", types.NewNumber(nil), ast.NewNumber(42))},
							ast.NewVariable("f"),
						),
					),
				},
			),
			ast.NewModule(
				"",
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewLambda([]string{"x"}, ast.NewVariable("x")),
					),
					ast.NewBind(
						"g",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewLambda(
							[]string{"additional.argument-0"},
							ast.NewLet(
								[]ast.Bind{ast.NewBind("a", types.NewNumber(nil), ast.NewNumber(42))},
								ast.NewApplication(
									ast.NewVariable("f"),
									[]ast.Expression{ast.NewVariable("additional.argument-0")},
								),
							),
						),
					),
				},
			),
		},
		// Binds in let expressions
		{
			ast.NewModule(
				"",
				[]ast.Bind{
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
									"g",
									types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
									ast.NewVariable("f"),
								),
							},
							ast.NewNumber(42),
						),
					),
				},
			),
			ast.NewModule(
				"",
				[]ast.Bind{
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
									"g",
									types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
									ast.NewLambda(
										[]string{"additional.argument-0"},
										ast.NewApplication(
											ast.NewVariable("f"),
											[]ast.Expression{ast.NewVariable("additional.argument-0")},
										),
									),
								),
							},
							ast.NewNumber(42),
						),
					),
				},
			),
		},
		// Lambda expressions
		{
			ast.NewModule(
				"",
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewLambda([]string{"x"}, ast.NewVariable("x")),
					),
					ast.NewBind(
						"g",
						types.NewFunction(
							types.NewNumber(nil),
							types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
							nil,
						),
						ast.NewLambda([]string{"x"}, ast.NewVariable("f")),
					),
				},
			),
			ast.NewModule(
				"",
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewLambda([]string{"x"}, ast.NewVariable("x")),
					),
					ast.NewBind(
						"g",
						types.NewFunction(
							types.NewNumber(nil),
							types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
							nil,
						),
						ast.NewLambda(
							[]string{"x", "additional.argument-0"},
							ast.NewApplication(
								ast.NewVariable("f"),
								[]ast.Expression{ast.NewVariable("additional.argument-0")},
							),
						),
					),
				},
			),
		},
	} {
		assert.Equal(t, ms[1], desugarPartialApplications(ms[0]))
	}
}
