package command

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ein-lang/ein/command/compile"
	"github.com/ein-lang/ein/command/debug"
	"github.com/ein-lang/ein/command/generate"
	"github.com/ein-lang/ein/command/parse"
	"github.com/spf13/cobra"
)

var buildCommand = cobra.Command{
	Use:   "build <filename>",
	Short: "Build a source file into a binary",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, as []string) {
		if err := runBuildCommand(as[0]); err != nil {
			fmt.Fprintln(os.Stderr, err)

			if err, ok := err.(debug.Error); ok {
				fmt.Fprintln(os.Stderr, err.DebugInformation())
			}

			os.Exit(1)
		}
	},
}

func runBuildCommand(f string) error {
	bs, err := ioutil.ReadFile(f)

	if err != nil {
		return err
	}

	m, err := parse.Parse(f, string(bs))

	if err != nil {
		return err
	}

	r, err := getRuntimePath()

	if err != nil {
		return err
	}

	mm, err := compile.Compile(m)

	if err != nil {
		return err
	}

	return generate.Executable(mm, f, r)
}
