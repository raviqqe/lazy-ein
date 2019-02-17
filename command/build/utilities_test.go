package build

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setUpEnvironmentDirectories(t *testing.T) (string, string, func()) {
	cacheDir, err := ioutil.TempDir("", "")
	assert.Nil(t, err)

	rootDir, err := ioutil.TempDir("", "")
	assert.Nil(t, err)

	return cacheDir, rootDir, func() {
		os.Remove(cacheDir)
		os.Remove(rootDir)
	}
}

func TestNormalizePath(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	assert.Nil(t, err)

	f, err := ioutil.TempFile(d, "")
	assert.Nil(t, err)

	s, err := normalizePath(f.Name(), d)
	assert.Nil(t, err)
	assert.Equal(t, filepath.Base(f.Name()), s)
}
