package desugar

import (
	"testing"

	"github.com/ein-lang/ein/command/ast"
	"github.com/stretchr/testify/assert"
)

func TestDesugar(t *testing.T) {
	for _, ms := range [][2]ast.Module{
		{
			ast.NewModule("", []ast.Bind{}),
			ast.NewModule("", []ast.Bind{}),
		},
		// TODO: Add more test cases to check interation of desugar functions.
	} {
		assert.Equal(t, ms[1], Desugar(ms[0]))
	}
}
