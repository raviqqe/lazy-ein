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
	f, err := ioutil.TempFile("", "")
	defer os.Remove(f.Name())
	assert.Nil(t, err)

	_, err = f.Write([]byte(source))
	assert.Nil(t, err)

	d, err := ioutil.TempDir("", "")
	defer os.Remove(d)
	assert.Nil(t, err)

	assert.Nil(t, build.Executable(f.Name(), "../..", d))
}

func TestExecutableWithoutMainFunction(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	defer os.Remove(f.Name())
	assert.Nil(t, err)

	d, err := ioutil.TempDir("", "")
	defer os.Remove(d)
	assert.Nil(t, err)

	assert.Equal(
		t,
		errors.New("main function not found"),
		build.Executable(f.Name(), "../..", d),
	)
}
