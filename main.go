package main

import (
	"log"
	"foundry/cli/cmd"
	// "foundry/cli/auth"
	// "context"
)

func init() {
	// Remove timestamp prefix
	log.SetFlags(0)
}

func main() {
	cmd.Execute()

	// a := auth.New()
	// a.SignIn(context.TODO(), "vasek@foundryapp.co", "123456")

	// log.Println(a)
}
