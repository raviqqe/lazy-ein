package build_test

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ein-lang/ein/command/build"
	"github.com/stretchr/testify/assert"
)

const source = "main : Number -> Number\nmain x = 42"

func TestExecutable(t *testing.T) {
	cacheDir, err := ioutil.TempDir("", "")
	defer os.Remove(cacheDir)
	assert.Nil(t, err)

	rootDir, err := ioutil.TempDir("", "")
	defer os.Remove(cacheDir)
	assert.Nil(t, err)

	f, err := ioutil.TempFile(rootDir, "")
	defer os.Remove(f.Name())
	assert.Nil(t, err)

	_, err = f.WriteString(source)
	assert.Nil(t, err)

	assert.Nil(t, build.Executable(f.Name(), "../..", rootDir, cacheDir))
}

func TestExecutableWithoutMainFunction(t *testing.T) {
	cacheDir, err := ioutil.TempDir("", "")
	defer os.Remove(cacheDir)
	assert.Nil(t, err)

	rootDir, err := ioutil.TempDir("", "")
	defer os.Remove(cacheDir)
	assert.Nil(t, err)

	f, err := ioutil.TempFile(rootDir, "")
	defer os.Remove(f.Name())
	assert.Nil(t, err)

	assert.Equal(
		t,
		errors.New("main function not found"),
		build.Executable(f.Name(), "../..", rootDir, cacheDir),
	)
}
