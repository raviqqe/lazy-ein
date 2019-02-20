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

	n := filepath.Join(rootDir, "main.ein")
	assert.Nil(t, ioutil.WriteFile(n, []byte(source), 0644))
	defer os.Remove(n)

	assert.Nil(t, Build(n, "../..", rootDir, cacheDir))

	_, err := os.Stat("a.out")
	assert.Nil(t, err)

	os.Remove("a.out")
}

func TestBuildWithSubmodules(t *testing.T) {
	cacheDir, rootDir, clean := setUpEnvironmentDirectories(t)
	defer clean()

	n := filepath.Join(rootDir, "foo.ein")
	assert.Nil(t, ioutil.WriteFile(n, nil, 0644))
	defer os.Remove(n)

	assert.Nil(t, Build(n, "../..", rootDir, cacheDir))

	_, err := os.Stat("a.out")
	assert.Error(t, err)
}

func TestBuildWithMainModulesImportingOtherModules(t *testing.T) {
	cacheDir, rootDir, clean := setUpEnvironmentDirectories(t)
	defer clean()

	n := filepath.Join(rootDir, "foo.ein")
	assert.Nil(t, ioutil.WriteFile(n, []byte("export { y }\ny : Number\ny = 42"), 0644))
	defer os.Remove(n)

	n = filepath.Join(rootDir, "main.ein")
	assert.Nil(
		t,
		ioutil.WriteFile(n, []byte("import \"foo\"\nmain : Number -> Number\nmain x = foo.y"), 0644),
	)
	defer os.Remove(n)

	assert.Nil(t, Build(n, "../..", rootDir, cacheDir))

	_, err := os.Stat("a.out")
	assert.Nil(t, err)

	os.Remove("a.out")
}
