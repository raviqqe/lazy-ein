package command

import (
	"io/ioutil"

	"github.com/raviqqe/jsonxx/command/compile"
	"github.com/raviqqe/jsonxx/command/generate"
	"github.com/raviqqe/jsonxx/command/parse"
	"github.com/spf13/cobra"
)

var buildCommand = cobra.Command{
	Use:   "build",
	Short: "Build a source file into a binary",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, as []string) error {
		bs, err := ioutil.ReadFile(as[0])

		if err != nil {
			return err
		}

		m, err := parse.Parse(as[0], string(bs))

		if err != nil {
			return err
		}

		return generate.Executable(as[0], compile.Compile(m))
	},
}
