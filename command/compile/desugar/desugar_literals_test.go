package desugar

import (
	"testing"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestDesugarLiterals(t *testing.T) {
	for _, ms := range [][2]ast.Module{
		// Don't convert empty modules
		{
			ast.NewModule("", []ast.Bind{}),
			ast.NewModule("", []ast.Bind{}),
		},
		// Convert variable binds
		// TODO: Don't convert variable binds in a special way and optimize codes in core language.
		{
			ast.NewModule(
				"",
				[]ast.Bind{ast.NewBind("x", nil, types.NewNumber(nil), ast.NewNumber(42))},
			),
			ast.NewModule(
				"",
				[]ast.Bind{
					ast.NewBind("x", nil, types.NewUnboxed(types.NewNumber(nil), nil), ast.NewNumber(42)),
				},
			),
		},
		// Convert function binds
		{
			ast.NewModule(
				"foo",
				[]ast.Bind{
					ast.NewBind(
						"f",
						[]string{"x"},
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewNumber(42)),
				},
			),
			ast.NewModule(
				"foo",
				[]ast.Bind{
					ast.NewBind(
						"foo.literal-0",
						nil,
						types.NewUnboxed(types.NewNumber(nil), nil),
						ast.NewNumber(42),
					),
					ast.NewBind(
						"f",
						[]string{"x"},
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewVariable("foo.literal-0")),
				},
			),
		},
		// Don't convert non-literal expressions
		{
			ast.NewModule(
				"",
				[]ast.Bind{
					ast.NewBind(
						"f",
						[]string{"x"},
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewVariable("x")),
				},
			),
			ast.NewModule(
				"",
				[]ast.Bind{
					ast.NewBind(
						"f",
						[]string{"x"},
						types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
						ast.NewVariable("x")),
				},
			),
		},
	} {
		assert.Equal(t, ms[1], desugarLiterals(ms[0]))
	}
}
