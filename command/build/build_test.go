package build

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const source = "main : Number -> Number\nmain x = 42"

func TestBuildWithMainModules(t *testing.T) {
	cacheDir, rootDir, clean := setUpEnvironmentDirectories(t)
	defer clean()

	f, err := ioutil.TempFile(rootDir, "")
	defer os.Remove(f.Name())
	assert.Nil(t, err)

	_, err = f.WriteString(source)
	assert.Nil(t, err)

	assert.Nil(t, Build(f.Name(), "../..", rootDir, cacheDir))

	_, err = os.Stat("a.out")
	assert.Nil(t, err)

	os.Remove("a.out")
}

func TestBuildWithSubmodules(t *testing.T) {
	cacheDir, rootDir, clean := setUpEnvironmentDirectories(t)
	defer clean()

	f, err := ioutil.TempFile(rootDir, "")
	defer os.Remove(f.Name())
	assert.Nil(t, err)

	assert.Nil(t, Build(f.Name(), "../..", rootDir, cacheDir))

	_, err = os.Stat("a.out")
	assert.Error(t, err)
}
