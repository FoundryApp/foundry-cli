package cmd

import (
	"os"
	"foundry/cli/prompt"
)

func Exit() *prompt.Cmd {
	return &prompt.Cmd{"exit", "Stop Foundry CLI", runExit}
}

func runExit(args []string) error {
	os.Exit(0)
	return nil
}