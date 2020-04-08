package cmd

import (
	"fmt"
	"foundry/cli/firebase"
	"foundry/cli/logger"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	envDelCmd = &cobra.Command{
		Use:     "env-delete",
		Short:   "Delete environment variable(s) from your cloud environment",
		Example: "foundry env-delete ENV_1 ENV_2",
		Run:     runEnvDel,
	}
)

func init() {
	rootCmd.AddCommand(envDelCmd)
}

func runEnvDel(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		logger.WarningLogln("No envs to delete specified. Example usage: 'foundry env-delete ENV_1 ENV_2'")
		os.Exit(0)
	}

	reqBody := struct {
		Delete []string `json:"delete"`
	}{args}

	s := fmt.Sprintf("Will delete following env variables '%s'...", strings.Join(args, ","))
	logger.Logln(s)

	res, err := firebase.Call("deleteUserEnvs", authClient.IDToken, reqBody)
	if err != nil {
		logger.FdebuglnFatal("Error calling deleteUserEnvs:", err)
		logger.FatalLogln("Error deleting environment variables (1):", err)
	}
	if res.Error != nil {
		logger.FdebuglnFatal("Error calling deleteUserEnvs:", res.Error)
		logger.FatalLogln("Error deleting environment variables (2):", res.Error)
	}

	logger.SuccessLogln("Env variables deleted")

}
