package cmd

import "github.com/spf13/cobra"

var (
	getEnvsCmd = &cobra.Command{
		Use:   "get-envs",
		Short: "Prints all environment variables in your development environment",
		Long:  "",
		Run:   runGetEnvs,
	}
)

func init() {
	rootCmd.AddCommand(getEnvsCmd)
}

func runGetEnvs(cmd *cobra.Command, args []string) {
	// TODO: We should have a cloud func
}
