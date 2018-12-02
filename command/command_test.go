package command_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/raviqqe/jsonxx/command"
	"github.com/stretchr/testify/assert"
)

func TestCommandWithEmptySource(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	assert.Nil(t, err)
	defer os.Remove(f.Name())

	assert.Nil(t, command.Command([]string{"foo", f.Name()}))
}
