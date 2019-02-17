package command_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/ein-lang/ein/command"
	"github.com/stretchr/testify/assert"
)

func TestRootCommandError(t *testing.T) {
	command.Command.SetArgs(nil)
	assert.Error(t, command.Command.Execute())
}

func TestBuildCommand(t *testing.T) {
	os.Setenv("EIN_RUNTIME_PATH", "..")
	os.Setenv("EIN_MODULE_ROOT_PATH", ".")

	f, err := ioutil.TempFile("", "")
	assert.Nil(t, err)
	f.WriteString("main : Number -> Number\nmain x = 42")
	defer os.Remove(f.Name())

	command.Command.SetArgs([]string{"build", f.Name()})
	assert.Nil(t, command.Command.Execute())
}

func TestBuildCommandErrorWithInvalidArgument(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	assert.Nil(t, err)
	defer os.Remove(f.Name())

	command.Command.SetArgs([]string{"build", "--invalid-option", f.Name()})
	assert.Error(t, command.Command.Execute())
}
