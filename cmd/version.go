package cmd

import (
	"foundry/cli/logger"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:     "version",
		Short:   "Print version of Foundry",
		Example: "foundry version",
		Run:     runVersion,
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	logger.Logln("Foundry v0.3.0\n")
}
