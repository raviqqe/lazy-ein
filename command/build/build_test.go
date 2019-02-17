package build

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const source = "main : Number -> Number\nmain x = 42"

func TestBuild(t *testing.T) {
	cacheDir, rootDir, clean := setUpEnvironmentDirectories(t)
	defer clean()

	f, err := ioutil.TempFile(rootDir, "")
	defer os.Remove(f.Name())
	assert.Nil(t, err)

	_, err = f.WriteString(source)
	assert.Nil(t, err)

	assert.Nil(t, Build(f.Name(), "../..", rootDir, cacheDir))
}

func TestBuildWithoutMainFunction(t *testing.T) {
	cacheDir, rootDir, clean := setUpEnvironmentDirectories(t)
	defer clean()

	f, err := ioutil.TempFile(rootDir, "")
	defer os.Remove(f.Name())
	assert.Nil(t, err)

	assert.Equal(
		t,
		errors.New("main function not found"),
		Build(f.Name(), "../..", rootDir, cacheDir),
	)
}
