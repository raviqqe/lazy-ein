package generate_test

import (
	"testing"

	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/generate"
	"github.com/stretchr/testify/assert"
)

func TestExecutable(t *testing.T) {
	assert.Nil(
		t,
		generate.Executable(
			ast.NewModule(
				"main.ein",
				nil,
				[]ast.Bind{
					ast.NewBind(
						"main",
						ast.NewLambda(
							nil,
							false,
							[]ast.Argument{ast.NewArgument("x", types.NewBoxed(types.NewFloat64()))},
							ast.NewApplication(ast.NewVariable("x"), nil),
							types.NewBoxed(types.NewFloat64()),
						),
					),
				},
			),
			"main.ein",
			"../..",
		),
	)
}

func TestExecutableErrorWithoutMainFunction(t *testing.T) {
	assert.Error(
		t,
		generate.Executable(ast.NewModule("main.ein", nil, nil), "main.ein", "../.."),
	)
}
