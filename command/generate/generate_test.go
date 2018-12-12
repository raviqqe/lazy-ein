package generate_test

import (
	"testing"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/compile"
	"github.com/ein-lang/ein/command/generate"
	"github.com/ein-lang/ein/command/types"
	"github.com/stretchr/testify/assert"
	"llvm.org/llvm/bindings/go/llvm"
)

func TestExecutable(t *testing.T) {
	m, err := compile.Compile(ast.NewModule(
		"main.ein",
		[]ast.Bind{
			ast.NewBind(
				"main",
				[]string{"x"},
				types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
				ast.NewNumber(42),
			),
		},
	))

	assert.Nil(t, err)

	assert.Nil(
		t,
		generate.Executable(
			m,
			"main.ein",
			"../..",
		),
	)
}

func TestExecutableErrorWithoutMainFunction(t *testing.T) {
	assert.Error(
		t,
		generate.Executable(llvm.NewModule("main.ein"), "main.ein", "../.."),
	)
}
