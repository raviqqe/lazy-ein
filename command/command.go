package command

import (
	"io/ioutil"

	"github.com/raviqqe/jsonxx/command/compile"
	"github.com/raviqqe/jsonxx/command/generate"
	"github.com/raviqqe/jsonxx/command/parse"
)

// Command runs a compiler command with command-line arguments.
func Command(ss []string) error {
	as, err := getArguments(ss)

	if err != nil {
		return err
	}

	bs, err := ioutil.ReadFile(as.Filename)

	if err != nil {
		return err
	}

	m, err := parse.Parse(as.Filename, string(bs))

	if err != nil {
		return err
	}

	return generate.Executable(as.Filename, compile.Compile(m))
}
