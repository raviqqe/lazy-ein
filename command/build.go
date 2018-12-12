package command

import (
	"io/ioutil"

	"github.com/ein-lang/ein/command/compile"
	"github.com/ein-lang/ein/command/generate"
	"github.com/ein-lang/ein/command/parse"
	"github.com/spf13/cobra"
)

var buildCommand = cobra.Command{
	Use:   "build <filename>",
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

		r, err := getRuntimeRoot()

		if err != nil {
			return err
		}

		return generate.Executable(compile.Compile(m), as[0], r)
	},
}
