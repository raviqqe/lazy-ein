package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewModuleName(t *testing.T) {
	n, err := NewModuleName("/foo/bar", "/foo")
	assert.Nil(t, err)
	assert.Equal(t, ModuleName("bar"), n)
}

func TestNewModuleNameWithModulesOutsideOfRootDirectories(t *testing.T) {
	n, err := NewModuleName("/bar/baz", "/foo")
	assert.Nil(t, err)
	assert.Equal(t, ModuleName("/bar/baz"), n)
}

func TestModuleNameToPath(t *testing.T) {
	n, err := NewModuleName("/foo/bar", "/foo")
	assert.Nil(t, err)

	assert.Equal(t, "/foo/bar", n.ToPath("/foo"))
}
