package desugar

import (
	"testing"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestDesugarLists(t *testing.T) {
	for _, ms := range [][2]ast.Module{
		// Empty modules
		{
			ast.NewModule(ast.NewExport(), nil, []ast.Bind{}),
			ast.NewModule(ast.NewExport(), nil, []ast.Bind{}),
		},
		// Simple lists
		{
			ast.NewModule(
				ast.NewExport(),
				nil,
				[]ast.Bind{
					ast.NewBind(
						"x",
						types.NewList(types.NewNumber(nil), nil),
						ast.NewList(
							types.NewList(types.NewNumber(nil), nil),
							[]ast.ListArgument{ast.NewListArgument(ast.NewVariable("y"), false)},
						),
					),
				},
			),
			ast.NewModule(
				ast.NewExport(),
				nil,
				[]ast.Bind{
					ast.NewBind(
						"x",
						types.NewList(types.NewNumber(nil), nil),
						ast.NewList(
							types.NewList(types.NewNumber(nil), nil),
							[]ast.ListArgument{ast.NewListArgument(ast.NewVariable("y"), false)},
						),
					),
				},
			),
		},
		// Complex lists
		{
			ast.NewModule(
				ast.NewExport(),
				nil,
				[]ast.Bind{
					ast.NewBind(
						"x",
						types.NewList(types.NewNumber(nil), nil),
						ast.NewList(
							types.NewList(types.NewNumber(nil), nil),
							[]ast.ListArgument{
								ast.NewListArgument(
									ast.NewApplication(
										ast.NewVariable("f"),
										[]ast.Expression{ast.NewVariable("y")},
									),
									false,
								),
							},
						),
					),
				},
			),
			ast.NewModule(
				ast.NewExport(),
				nil,
				[]ast.Bind{
					ast.NewBind(
						"x",
						types.NewList(types.NewNumber(nil), nil),
						ast.NewLet(
							[]ast.Bind{
								ast.NewBind(
									"$list.element-0",
									types.NewUnknown(nil),
									ast.NewApplication(
										ast.NewVariable("f"),
										[]ast.Expression{ast.NewVariable("y")},
									),
								),
							},
							ast.NewList(
								types.NewList(types.NewNumber(nil), nil),
								[]ast.ListArgument{
									ast.NewListArgument(ast.NewVariable("$list.element-0"), false),
								},
							),
						),
					),
				},
			),
		},
	} {
		assert.Equal(t, ms[1], desugarLists(ms[0]))
	}
}
