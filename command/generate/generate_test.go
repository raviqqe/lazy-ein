package generate_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/raviqqe/jsonxx/command/core/ast"
	"github.com/raviqqe/jsonxx/command/generate"
	"github.com/stretchr/testify/assert"
)

func TestModule(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	assert.Nil(t, err)
	defer os.Remove(f.Name())

	assert.Nil(
		t,
		generate.Executable(
			f.Name(),
			ast.NewModule(f.Name(), nil, nil),
		),
	)
}
