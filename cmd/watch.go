package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
  Use:   "",
  Short: "",
  Long:  ``,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Foundry watch")
  },
}