package parse_test

import (
	"testing"

	"github.com/raviqqe/jsonxx/ast"
	"github.com/raviqqe/jsonxx/parse"
	"github.com/stretchr/testify/assert"
)

func TestParseWithEmptySource(t *testing.T) {
	x, err := parse.Parse("", "")

	assert.Equal(t, []ast.Bind{}, x)
	assert.Nil(t, err)
}

func TestParseWithVariableBind(t *testing.T) {
	x, err := parse.Parse("", "x = 42")

	assert.Equal(t, []ast.Bind{ast.NewBind("x", ast.NewNumber(42))}, x)
	assert.Nil(t, err)
}

func TestParseWithVariousIdentifiers(t *testing.T) {
	for _, s := range []string{
		"az = 42",
		"a09 = 42",
	} {
		_, err := parse.Parse("", s)

		assert.Nil(t, err)
	}
}

func TestParseWithInvalidIdentifiers(t *testing.T) {
	for _, s := range []string{
		" = 42",
		"0 = 42",
		"1x = 42",
		"let = 42",
	} {
		_, err := parse.Parse("", s)

		assert.NotNil(t, err)
	}
}

func TestParseWithNumbers(t *testing.T) {
	for _, s := range []string{
		"x = 1",
		"x = 42",
		"x = 1.0",
		"x = -1",
		"x = 0",
	} {
		_, err := parse.Parse("", s)

		assert.Nil(t, err)
	}
}

func TestParseWithInvalidNumbers(t *testing.T) {
	for _, s := range []string{
		"x = - 42",
		"x = 01",
		"x = 1.",
	} {
		_, err := parse.Parse("", s)

		assert.NotNil(t, err)
	}
}
