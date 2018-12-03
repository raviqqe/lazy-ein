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
			[]ast.Bind{ast.NewBind("x", types.NewNumber(debugInformation), ast.NewNumber(42))},
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
