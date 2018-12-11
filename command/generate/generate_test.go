package generate_test

import (
	"os"
	"testing"

	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/generate"
	"github.com/stretchr/testify/assert"
)

func TestExecutable(t *testing.T) {
	os.Setenv("EIN_ROOT", "../..")

	assert.Nil(
		t,
		generate.Executable(
			"main.ein",
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
		),
	)
}

func TestExecutableErrorWithoutMainFunction(t *testing.T) {
	os.Setenv("EIN_ROOT", "../..")

	assert.Error(
		t,
		generate.Executable("main.ein", ast.NewModule("main.ein", nil, nil)),
	)
}

func TestExecutableErrorWithoutRootEnvironmentVariable(t *testing.T) {
	os.Unsetenv("EIN_ROOT")

	assert.Error(
		t,
		generate.Executable(
			"main.ein",
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
		),
	)
}
