package command

import (
	"errors"

	"github.com/spf13/cobra"
)

// Command is a command.
var Command = func() cobra.Command {
	c := cobra.Command{
		Use: "ein",
		RunE: func(*cobra.Command, []string) error {
			return errors.New("a subcommand must be provided")
		},
		Short:   "Ein programming language",
		Version: "0.0.0",
	}

	c.AddCommand(&buildCommand)

	return c
}()
