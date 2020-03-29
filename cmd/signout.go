package cmd

import (
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
	authClient.ClearTokens()
	color.Green("✔ Signed Out")
}