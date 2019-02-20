package build

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

func TestBuildWithMainModulesImportingOtherModules(t *testing.T) {
	cacheDir, rootDir, clean := setUpEnvironmentDirectories(t)
	defer clean()

	n := filepath.Join(rootDir, "bar")
	assert.Nil(t, ioutil.WriteFile(n, []byte("export { y }\ny : Number\ny = 42"), 0644))
	defer os.Remove(n)

	f, err := ioutil.TempFile(rootDir, "")
	defer os.Remove(f.Name())
	assert.Nil(t, err)

	_, err = f.WriteString("import \"bar\"\nmain : Number -> Number\nmain x = bar.y")
	assert.Nil(t, err)

	assert.Nil(t, Build(f.Name(), "../..", rootDir, cacheDir))

	_, err = os.Stat("a.out")
	assert.Nil(t, err)

	os.Remove("a.out")
}
