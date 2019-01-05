package compile

import (
	"testing"

	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/types"
	"github.com/stretchr/testify/assert"
	"llvm.org/llvm/bindings/go/llvm"
)

func TestNewTypeGenerator(t *testing.T) {
	_, err := newTypeGenerator(
		llvm.NewModule("foo"),
		[]ast.TypeDefinition{ast.NewTypeDefinition("foo", types.NewFloat64())},
	)
	assert.Nil(t, err)
}

func TestNewTypeGeneratorWithRecursiveTypes(t *testing.T) {
	_, err := newTypeGenerator(
		llvm.NewModule("foo"),
		[]ast.TypeDefinition{
			ast.NewTypeDefinition(
				"foo",
				types.NewAlgebraic(
					[]types.Constructor{
						types.NewConstructor(
							"constructor",
							[]types.Type{types.NewBoxed(types.NewNamed("foo"))},
						),
					},
				),
			),
		},
	)
	assert.Nil(t, err)
}

func TestTypeGeneratorGenerateWithNamedTypes(t *testing.T) {
	g, err := newTypeGenerator(
		llvm.NewModule("foo"),
		[]ast.TypeDefinition{
			ast.NewTypeDefinition(
				"foo",
				types.NewAlgebraic(
					[]types.Constructor{
						types.NewConstructor("constructor", []types.Type{types.NewFloat64()}),
					},
				),
			),
		},
	)
	assert.Nil(t, err)

	g.Generate(types.NewNamed("foo"))
}

func TestTypeGeneratorGenerateErrorWithUndefinedTypes(t *testing.T) {
	g, err := newTypeGenerator(llvm.NewModule("foo"), nil)
	assert.Nil(t, err)

	_, err = g.Generate(types.NewNamed("foo"))
	assert.Error(t, err)
}

func TestTypeGeneratorGenerateSizedPayload(t *testing.T) {
	for _, c := range []struct {
		lambda ast.Lambda
		size   int
	}{
		{
			lambda: ast.NewLambda(nil, true, nil, ast.NewFloat64(42), types.NewFloat64()),
			size:   8,
		},
		{
			lambda: ast.NewLambda(
				[]ast.Argument{
					ast.NewArgument("x", types.NewFloat64()),
					ast.NewArgument("y", types.NewFloat64()),
				},
				true,
				nil,
				ast.NewPrimitiveOperation(
					ast.AddFloat64,
					[]ast.Atom{ast.NewVariable("x"), ast.NewVariable("y")},
				),
				types.NewFloat64(),
			),
			size: 16,
		},
		{
			lambda: ast.NewLambda(nil, false, nil, ast.NewFloat64(42), types.NewFloat64()),
			size:   0,
		},
	} {
		g, err := newTypeGenerator(llvm.NewModule("foo"), nil)
		assert.Nil(t, err)

		tt, err := g.generateSizedPayload(c.lambda)
		assert.Nil(t, err)

		assert.Equal(t, c.size, tt.ArrayLength())
	}
}

func TestTypeGeneratorBytesToWords(t *testing.T) {
	g, err := newTypeGenerator(llvm.NewModule("foo"), nil)
	assert.Nil(t, err)

	for k, v := range map[int]int{
		0:  0,
		1:  1,
		2:  1,
		7:  1,
		8:  1,
		9:  2,
		16: 2,
		17: 3,
	} {
		assert.Equal(t, v, g.bytesToWords(k))
	}
}
