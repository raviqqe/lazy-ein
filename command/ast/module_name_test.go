package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewModuleName(t *testing.T) {
	s, err := NewModuleName("/foo/bar", "/foo")
	assert.Nil(t, err)
	assert.Equal(t, ModuleName("bar"), s)
}

func TestNewModuleNameWithModulesOutsideOfRootDirectories(t *testing.T) {
	s, err := NewModuleName("/bar/baz", "/foo")
	assert.Nil(t, err)
	assert.Equal(t, ModuleName("/bar/baz"), s)
}
