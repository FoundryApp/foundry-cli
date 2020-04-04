package main

import (
	"foundry/cli/cmd"
	"foundry/cli/config"
	"foundry/cli/logger"
)

func main() {
	if err := config.Init(); err != nil {
		logger.FatalLogln("Couldn't init config", err)
	}

	cmd.Execute()
}
