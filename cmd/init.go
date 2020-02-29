package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(initCmd)
}





var initCmd = &cobra.Command{
  Use:   "init",
  Short: "Init Foundry",
  Long:  ``,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
  },
}