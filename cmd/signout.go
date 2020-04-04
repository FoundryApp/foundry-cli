package cmd

import (
	"foundry/cli/logger"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	signOutCmd = &cobra.Command{
		Use:   "sign-out",
		Short: "Sign out",
		Long:  "",
		Run:   runSignOut,
	}
)

func init() {
	rootCmd.AddCommand(signOutCmd)
}

func runSignOut(cmd *cobra.Command, args []string) {
	if err := authClient.SignOut(); err != nil {
		logger.FdebuglnFatal("Sign out error", err)
		logger.FatalLogln("Sign out error", err)
	}
	color.Green("✔ Signed Out")
}
