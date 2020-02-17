package main

import (
	"fmt"
	"os"

	"github.com/raviqqe/lazy-ein/command"
)

func main() {
	if err := command.Command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
