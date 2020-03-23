package cmd

import (
	"fmt"
	"foundry/cli/prompt"
)

func Watch() *prompt.Cmd {
	return &prompt.Cmd{"watch", "Watch only specific function(s)", runWatch}
}

func runWatch(args []string) error {
	fmt.Println("watching", args)
	return nil
}