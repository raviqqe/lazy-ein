package desugar

import (
	"testing"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestDesugarComplexLists(t *testing.T) {
	for _, ms := range [][2]ast.Module{
		// Empty modules
		{
			ast.NewModule("", []ast.Bind{}),
			ast.NewModule("", []ast.Bind{}),
		},
		// Simple lists
		{
			ast.NewModule(
				"foo",
				[]ast.Bind{
					ast.NewBind(
						"x",
						types.NewList(types.NewNumber(nil), nil),
						ast.NewList(
							types.NewList(types.NewNumber(nil), nil),
							[]ast.Expression{ast.NewVariable("y")}),
					),
				},
			),
			ast.NewModule(
				"foo",
				[]ast.Bind{
					ast.NewBind(
						"x",
						types.NewList(types.NewNumber(nil), nil),
						ast.NewList(
							types.NewList(types.NewNumber(nil), nil),
							[]ast.Expression{ast.NewVariable("y")}),
					),
				},
			),
		},
		// Complex lists
		{
			ast.NewModule(
				"foo",
				[]ast.Bind{
					ast.NewBind(
						"x",
						types.NewList(types.NewNumber(nil), nil),
						ast.NewList(
							types.NewList(types.NewNumber(nil), nil),
							[]ast.Expression{
								ast.NewApplication(ast.NewVariable("f"), []ast.Expression{ast.NewVariable("y")}),
							},
						),
					),
				},
			),
			ast.NewModule(
				"foo",
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
								[]ast.Expression{ast.NewVariable("$list.element-0")},
							),
						),
					),
				},
			),
		},
	} {
		assert.Equal(t, ms[1], desugarComplexLists(ms[0]))
	}
}
