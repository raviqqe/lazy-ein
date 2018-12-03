package parse

import (
	"testing"

	"github.com/raviqqe/jsonxx/command/ast"
	"github.com/raviqqe/jsonxx/command/debug"
	"github.com/raviqqe/jsonxx/command/types"
	"github.com/stretchr/testify/assert"
)

func TestParseWithEmptySource(t *testing.T) {
	x, err := Parse("", "")

	assert.Equal(t, ast.NewModule("", []ast.Bind{}), x)
	assert.Nil(t, err)
}

func TestParseWithVariableBind(t *testing.T) {
	x, err := Parse("", "x : Number\nx = 42\n")

	assert.Equal(
		t,
		ast.NewModule(
			"",
			[]ast.Bind{
				ast.NewBind(
					"x",
					types.NewNumber(debug.NewInformation("", 1, 5, "x : Number")),
					ast.NewNumber(42),
				),
			},
		),
		x,
	)
	assert.Nil(t, err)
}

func TestStateIdentifier(t *testing.T) {
	for _, s := range []string{"x", "az", "a09"} {
		_, err := newState("", s).identifier()()
		assert.Nil(t, err)
	}
}

func TestStateIdentifierError(t *testing.T) {
	for _, s := range []string{"0", "1x", "let"} {
		_, err := newState("", s).identifier()()
		assert.NotNil(t, err)
	}
}

func TestStateNumberLiteral(t *testing.T) {
	for _, s := range []string{
		"1",
		"42",
		"1.0",
		"-1",
		"0",
	} {
		_, err := newState("", s).numberLiteral()()
		assert.Nil(t, err)
	}
}

func TestStateNumberLiteralError(t *testing.T) {
	for _, s := range []string{
		"- 42",
		"01",
		"1.",
	} {
		s := newState("", s)
		_, err := s.Exhaust(s.numberLiteral())()
		assert.NotNil(t, err)
	}
}
