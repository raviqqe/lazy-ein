package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewModuleName(t *testing.T) {
	n, err := NewModuleName("/foo/bar.ein", "/foo")
	assert.Nil(t, err)
	assert.Equal(t, ModuleName("bar"), n)
}

func TestNewModuleNameWithModulesOutsideOfRootDirectories(t *testing.T) {
	n, err := NewModuleName("/bar/baz.ein", "/foo")
	assert.Nil(t, err)
	assert.Equal(t, ModuleName("/bar/baz"), n)
}

func TestModuleNameToPath(t *testing.T) {
	n, err := NewModuleName("/foo/bar.ein", "/foo")
	assert.Nil(t, err)

	assert.Equal(t, "/foo/bar.ein", n.ToPath("/foo"))
}
