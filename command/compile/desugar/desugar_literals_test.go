package desugar

import (
	"testing"

	"github.com/raviqqe/lazy-ein/command/ast"
	"github.com/raviqqe/lazy-ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestDesugarLiterals(t *testing.T) {
	for _, ms := range [][2]ast.Module{
		// Don't convert empty modules
		{
			ast.NewModule("", ast.NewExport(), nil, []ast.Bind{}),
			ast.NewModule("", ast.NewExport(), nil, []ast.Bind{}),
		},
		// Convert variable binds
		// TODO: Don't convert variable binds in a special way and optimize codes in core language.
		{
			ast.NewModule(
				"",
				ast.NewExport(),
				nil,
				[]ast.Bind{ast.NewBind("x", types.NewNumber(nil), ast.NewNumber(42))},
			),
			ast.NewModule(
				"",
				ast.NewExport(),
				nil,
				[]ast.Bind{
					ast.NewBind(
						"x",
						types.NewUnboxed(types.NewNumber(nil), nil),
						ast.NewUnboxed(ast.NewNumber(42)),
					),
				},
			),
		},
		// Convert function binds
		{
			ast.NewModule(
				"",
				ast.NewExport(),
				nil,
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewLambda([]string{"x"}, ast.NewNumber(42))),
				},
			),
			ast.NewModule(
				"",
				ast.NewExport(),
				nil,
				[]ast.Bind{
					ast.NewBind(
						"$literal-0",
						types.NewUnboxed(types.NewNumber(nil), nil),
						ast.NewUnboxed(ast.NewNumber(42)),
					),
					ast.NewBind(
						"f",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewLambda([]string{"x"}, ast.NewVariable("$literal-0"))),
				},
			),
		},
		// Don't convert non-literal expressions
		{
			ast.NewModule(
				"",
				ast.NewExport(),
				nil,
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewLambda([]string{"x"}, ast.NewVariable("x")),
					),
				},
			),
			ast.NewModule(
				"",
				ast.NewExport(),
				nil,
				[]ast.Bind{
					ast.NewBind(
						"f",
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewLambda([]string{"x"}, ast.NewVariable("x")),
					),
				},
			),
		},
	} {
		assert.Equal(t, ms[1], desugarLiterals(ms[0]))
	}
}
