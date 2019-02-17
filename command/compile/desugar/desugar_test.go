package desugar_test

import (
	"testing"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/compile/desugar"
	"github.com/stretchr/testify/assert"
)

func TestWithoutTypes(t *testing.T) {
	for _, ms := range [][2]ast.Module{
		{
			ast.NewModule(ast.NewExport(), nil, []ast.Bind{}),
			ast.NewModule(ast.NewExport(), nil, []ast.Bind{}),
		},
		// TODO: Add more test cases to check interation of desugar functions.
	} {
		assert.Equal(t, ms[1], desugar.WithoutTypes(ms[0]))
	}
}
