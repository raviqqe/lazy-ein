package compile_test

import (
	"testing"

	"github.com/raviqqe/jsonxx/command/ast"
	"github.com/raviqqe/jsonxx/command/compile"
	cast "github.com/raviqqe/jsonxx/command/core/ast"
	ctypes "github.com/raviqqe/jsonxx/command/core/types"
	"github.com/raviqqe/jsonxx/command/types"
	"github.com/stretchr/testify/assert"
)

func TestCompileWithEmptySource(t *testing.T) {
	m := compile.Compile(ast.NewModule("", []ast.Bind{}))

	assert.Equal(t, cast.NewModule("", nil, []cast.Bind{}), m)
}

func TestCompileWithVariableBinds(t *testing.T) {
	m := compile.Compile(
		ast.NewModule(
			"",
			[]ast.Bind{ast.NewBind("x", nil, types.NewNumber(nil), ast.NewNumber(42))},
		),
	)

	assert.Equal(
		t,
		cast.NewModule(
			"",
			nil,
			[]cast.Bind{
				cast.NewBind(
					"x",
					cast.NewLambda(nil, true, nil, cast.NewFloat64(42), ctypes.NewFloat64()),
				),
			},
		),
		m,
	)
}

func TestCompileWithFunctionBinds(t *testing.T) {
	m := compile.Compile(
		ast.NewModule(
			"foo",
			[]ast.Bind{
				ast.NewBind(
					"f",
					[]string{"x", "y"},
					types.NewFunction(
						types.NewNumber(nil),
						types.NewFunction(
							types.NewNumber(nil),
							types.NewNumber(nil),
							nil,
						),
						nil,
					),
					ast.NewNumber(42),
				),
			},
		),
	)

	assert.Equal(
		t,
		cast.NewModule(
			"foo",
			nil,
			[]cast.Bind{
				cast.NewBind(
					"foo.literal-0",
					cast.NewLambda(nil, true, nil, cast.NewFloat64(42), ctypes.NewFloat64()),
				),
				cast.NewBind(
					"f",
					cast.NewLambda(
						nil,
						false,
						[]cast.Argument{
							cast.NewArgument("x", ctypes.NewBoxed(ctypes.NewFloat64())),
							cast.NewArgument("y", ctypes.NewBoxed(ctypes.NewFloat64())),
						},
						cast.NewApplication(cast.NewVariable("foo.literal-0"), nil),
						ctypes.NewBoxed(ctypes.NewFloat64()),
					),
				),
			},
		),
		m,
	)
}
