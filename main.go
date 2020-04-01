package main

import (
	"os"

	"foundry/cli/cmd"
	"foundry/cli/config"
	"foundry/cli/logger"
)

func main() {
	// time.Sleep(time.Second * 20)

	if err := config.Init(); err != nil {
		logger.Logln("Couldn't init config", err)
		os.Exit(1)
	}

	cmd.Execute()
}
