package build

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectCacheStore(t *testing.T) {
	cacheDir, rootDir, clean := setUpEnvironmentDirectories(t)
	defer clean()

	f, err := ioutil.TempFile(rootDir, "")
	defer os.Remove(f.Name())
	assert.Nil(t, err)

	_, err = f.WriteString("main : Number -> Number\nmain x = 42")
	assert.Nil(t, err)

	c := newObjectCache(cacheDir, rootDir)

	s, ok, err := c.Get(f.Name())

	assert.Nil(t, err)
	assert.False(t, ok)
	assert.Equal(t, "", s)

	s, err = c.Store(f.Name(), []byte("baz"))

	assert.Nil(t, err)
	assert.NotEqual(t, "", s)

	ss, ok, err := c.Get(f.Name())

	assert.Nil(t, err)
	assert.True(t, ok)
	assert.Equal(t, s, ss)
}

func TestObjectCacheGeneratePath(t *testing.T) {
	cacheDir, rootDir, clean := setUpEnvironmentDirectories(t)
	defer clean()

	f, err := ioutil.TempFile(rootDir, "")
	defer os.Remove(f.Name())
	assert.Nil(t, err)

	_, err = f.WriteString("export { x }\nx : Number\nx = 42")
	assert.Nil(t, err)

	ff, err := ioutil.TempFile(rootDir, "")
	defer os.Remove(f.Name())
	assert.Nil(t, err)

	_, err = ff.WriteString(
		fmt.Sprintf("import \"%s\"\nmain : Number -> Number\nmain x = 42", f.Name()),
	)
	assert.Nil(t, err)

	c := newObjectCache(cacheDir, rootDir)
	s, err := c.generatePath(ff.Name())

	assert.Nil(t, err)
	assert.NotEqual(t, "", s)

	_, err = f.WriteAt([]byte("export { y }\ny : Number\ny = 42"), 0)
	assert.Nil(t, err)

	ss, err := c.generatePath(ff.Name())

	assert.Nil(t, err)
	assert.NotEqual(t, s, ss)
}
