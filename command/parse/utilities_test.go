package parse

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizePath(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	assert.Nil(t, err)

	f, err := ioutil.TempFile(d, "")
	assert.Nil(t, err)

	s, err := normalizePath(f.Name(), d)
	assert.Nil(t, err)
	assert.Equal(t, filepath.Base(f.Name()), s)
}
