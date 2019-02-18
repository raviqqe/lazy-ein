package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewModuleName(t *testing.T) {
	s, err := NewModuleName("/dir/foo", "/dir")
	assert.Nil(t, err)
	assert.Equal(t, ModuleName("foo"), s)
}
