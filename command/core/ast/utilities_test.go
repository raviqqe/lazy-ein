package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceVariable(t *testing.T) {
	assert.Equal(t, "bar", replaceVariable("foo", map[string]string{"foo": "bar"}))
	assert.Equal(t, "foo", replaceVariable("foo", map[string]string{"bar": "baz"}))
}

func TestRemoveVariables(t *testing.T) {
	for _, ms := range [][2]map[string]string{
		{{}, {}},
		{{"foo": "bar"}, {}},
		{{"foo": "bar", "baz": "blah"}, {"baz": "blah"}},
	} {
		assert.Equal(t, ms[1], removeVariables(ms[0], "foo"))
	}
}
