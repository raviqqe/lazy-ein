package desugar

import (
	"testing"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestDesugarApplications(t *testing.T) {
	for _, ms := range [][2]ast.Module{
		// Don't convert empty modules
		{
			ast.NewModule("", []ast.Bind{}),
			ast.NewModule("", []ast.Bind{}),
		},
		// Convert arguments
		{
			ast.NewModule(
				"foo",
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewLambda([]string{"x"}, ast.NewVariable("x")),
					),
					ast.NewBind(
						"x",
						types.NewNumber(nil),
						ast.NewApplication(
							ast.NewVariable("f"),
							[]ast.Expression{
								ast.NewLet(
									[]ast.Bind{ast.NewBind("y", types.NewUnknown(nil), ast.NewNumber(42))},
									ast.NewVariable("y"),
								),
							},
						),
					),
				},
			),
			ast.NewModule(
				"foo",
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
									"foo.application.argument-0",
									types.NewUnknown(nil),
									ast.NewLet(
										[]ast.Bind{
											ast.NewBind("y", types.NewUnknown(nil), ast.NewNumber(42)),
										},
										ast.NewVariable("y"),
									),
								),
							},
							ast.NewApplication(
								ast.NewVariable("f"),
								[]ast.Expression{
									ast.NewVariable("foo.application.argument-0"),
								},
							),
						),
					),
				},
			),
		},
		// TODO: Add test to convert functions.
	} {
		assert.Equal(t, ms[1], desugarApplications(ms[0]))
	}
}
