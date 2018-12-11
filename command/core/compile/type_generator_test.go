package compile

import (
	"testing"

	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/types"
	"github.com/stretchr/testify/assert"
	"llvm.org/llvm/bindings/go/llvm"
)

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
		assert.Equal(
			t,
			c.size,
			newTypeGenerator(llvm.NewModule("foo")).generateSizedPayload(c.lambda).ArrayLength(),
		)
	}
}

func TestTypeGeneratorBytesToWords(t *testing.T) {
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
		assert.Equal(t, v, newTypeGenerator(llvm.NewModule("foo")).bytesToWords(k))
	}
}
