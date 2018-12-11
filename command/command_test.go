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
	os.Setenv("EIN_ROOT", "..")

	f, err := ioutil.TempFile("", "")
	assert.Nil(t, err)
	f.WriteString("main : Number -> Number\nmain x = 42")
	defer os.Remove(f.Name())

	command.Command.SetArgs([]string{"build", f.Name()})
	assert.Nil(t, command.Command.Execute())
}

func TestBuildCommandErrorWithInvalidFilename(t *testing.T) {
	command.Command.SetArgs([]string{"build", "invalid-filename"})
	assert.Error(t, command.Command.Execute())
}

func TestBuildCommandErrorWithInvalidArgument(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	assert.Nil(t, err)
	defer os.Remove(f.Name())

	command.Command.SetArgs([]string{"build", "--invalid-option", f.Name()})
	assert.Error(t, command.Command.Execute())
}
