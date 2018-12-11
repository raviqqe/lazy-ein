package main

import (
	"fmt"
	"os"

	"github.com/ein-lang/ein/command"
)

func main() {
	if err := command.Command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
