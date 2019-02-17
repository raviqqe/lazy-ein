package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ein-lang/ein/command/build"
	"github.com/ein-lang/ein/command/debug"
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
	runtime, err := getRuntimePath()

	if err != nil {
		return err
	}

	root, err := getModulesRootPath()

	if err != nil {
		return err
	}

	c, err := getCacheDirectory()

	if err != nil {
		return err
	}

	return build.Build(f, runtime, root, c)
}

func getCacheDirectory() (string, error) {
	d, err := os.UserCacheDir()

	if err != nil {
		return "", err
	}

	d = filepath.Join(d, "ein-lang")

	if err := os.MkdirAll(d, 0755); err != nil {
		return "", err
	}

	return d, nil
}
