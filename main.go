package main

import (
	"os"

	"foundry/cli/logger"
	"foundry/cli/cmd"
	"foundry/cli/config"
)

func main() {
	if err := config.Init(); err != nil {
		logger.Log("Couldn't init config", err)
		os.Exit(1)
	}

	cmd.Execute()
}
