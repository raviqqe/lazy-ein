package build

import (
	"io/ioutil"
	"os"
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
