package build

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
)

func TestObjectCacheStore(t *testing.T) {
	cacheDir, rootDir, clean := setUpEnvironmentDirectories(t)
	defer clean()

	n := filepath.Join(rootDir, "main.ein")
	assert.Nil(t, ioutil.WriteFile(n, []byte("main : Number -> Number\nmain x = 42"), 0644))
	defer os.Remove(n)

	c := newObjectCache(cacheDir, rootDir)

	m, ok, err := c.Get(n)

	assert.Nil(t, err)
	assert.False(t, ok)
	assert.True(t, m.IsNil())

	assert.Nil(t, c.Store(n, llvm.NewModule("main")))

	m, ok, err = c.Get(n)

	assert.Nil(t, err)
	assert.True(t, ok)
	assert.False(t, m.IsNil())
}

func TestObjectCacheGeneratePath(t *testing.T) {
	cacheDir, rootDir, clean := setUpEnvironmentDirectories(t)
	defer clean()

	n := filepath.Join(rootDir, "foo.ein")
	assert.Nil(t, ioutil.WriteFile(n, []byte("export { x }\nx : Number\nx = 42"), 0644))
	defer os.Remove(n)

	n = filepath.Join(rootDir, "main.ein")
	assert.Nil(
		t,
		ioutil.WriteFile(n, []byte("import \"foo\"\nmain : Number -> Number\nmain x = 42"), 0644),
	)
	defer os.Remove(n)

	c := newObjectCache(cacheDir, rootDir)
	s, err := c.generatePath(n)

	assert.Nil(t, err)
	assert.NotEqual(t, "", s)

	assert.Nil(t, ioutil.WriteFile(n, []byte("export { x }\nx : Number\nx = 123"), 0644))

	ss, err := c.generatePath(n)

	assert.Nil(t, err)
	assert.NotEqual(t, s, ss)
}

func TestObjectCacheGeneratePathWithUnnormalizedModulePaths(t *testing.T) {
	cacheDir, rootDir, clean := setUpEnvironmentDirectories(t)
	defer clean()

	n := filepath.Join(rootDir, "foo.ein")
	assert.Nil(t, ioutil.WriteFile(n, []byte("export { x }\nx : Number\nx = 42"), 0644))
	defer os.Remove(n)

	c := newObjectCache(cacheDir, rootDir)
	s, err := c.generatePath(n)
	assert.Nil(t, err)
	assert.NotEqual(t, "", s)

	assert.Nil(t, os.Chdir(filepath.Dir(n)))

	ss, err := c.generatePath(filepath.Base(n))
	assert.Nil(t, err)
	assert.Equal(t, s, ss)
}
