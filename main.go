package main

import (
	"log"

	"foundry/cli/cmd"
	"foundry/cli/config"
)

func init() {
	// Remove timestamp prefix
	log.SetFlags(0)
}

func main() {
	if err := config.Init(); err != nil {
		log.Fatal("Couldn't init config", err)
	}

	cmd.Execute()
}
