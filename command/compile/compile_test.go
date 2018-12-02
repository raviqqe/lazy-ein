package compile_test

import (
	"testing"

	"github.com/raviqqe/jsonxx/command/ast"
	"github.com/raviqqe/jsonxx/command/compile"
	core "github.com/raviqqe/jsonxx/command/core/ast"
	"github.com/raviqqe/jsonxx/command/core/types"
	"github.com/stretchr/testify/assert"
)

func TestCompileWithEmptySource(t *testing.T) {
	m := compile.Compile(ast.NewModule("", []ast.Bind{}))

	assert.Equal(t, core.NewModule("", nil, []core.Bind{}), m)
}

func TestCompileWithVariableBinds(t *testing.T) {
	m := compile.Compile(ast.NewModule("", []ast.Bind{ast.NewBind("x", ast.NewNumber(42))}))

	assert.Equal(
		t,
		core.NewModule(
			"",
			nil,
			[]core.Bind{
				core.NewBind(
					"x",
					core.NewLambda(nil, true, nil, core.NewFloat64(42), types.NewFloat64()),
				),
			},
		),
		m,
	)
}
