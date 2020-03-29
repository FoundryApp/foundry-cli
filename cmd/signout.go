package cmd

import (
	"foundry/cli/logger"

	"github.com/spf13/cobra"
	"github.com/fatih/color"
)

var (
	signOutCmd = &cobra.Command{
		Use: 		"sign-out",
		Short: 	"Sign out",
		Long: 	"",
		Run:		runSignOut,
	}
)

func init() {
	rootCmd.AddCommand(signOutCmd)
}

func runSignOut(cmd *cobra.Command, args []string) {
	if err := authClient.SignOut(); err != nil {
		logger.Fdebugln(err)
		logger.ErrorLoglnFatal(err)
	}
	color.Green("✔ Signed Out")
}