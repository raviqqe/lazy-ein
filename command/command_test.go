package command_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/raviqqe/jsonxx/command"
	"github.com/stretchr/testify/assert"
)

func TestRootCommandError(t *testing.T) {
	command.Command.SetArgs(nil)
	assert.Error(t, command.Command.Execute())
}

func TestBuildCommandWithEmptySource(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	assert.Nil(t, err)
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
