package cmd

import (
	"fmt"

	"foundry/cli/firebase"
	"foundry/cli/logger"

	"github.com/spf13/cobra"
)

var (
	envPrintCmd = &cobra.Command{
		Use:     "env-print",
		Short:   "Print all environment variables in your cloud environment",
		Example: "foundry env-print",
		Run:     runEnvPrint,
	}
)

func init() {
	rootCmd.AddCommand(envPrintCmd)
}

func runEnvPrint(cmd *cobra.Command, args []string) {
	resp, err := firebase.Call("UserEnvs", authClient.IDToken, nil)
	if err != nil || resp.Error != nil {
		logger.FdebuglnFatal("Error calling getUserEnvs:", err)
		logger.FatalLogln("Error printing environment variables (1):", err)
	}
	if resp.Error != nil {
		logger.FdebuglnFatal("Error calling getUserEnvs:", resp.Error)
		logger.FatalLogln("Error printing environment variables (2):", resp.Error)
	}

	envs, ok := resp.Result.(map[string]interface{})
	if !ok {
		logger.FdebuglnFatal("Failed to type assert res.Result")
		logger.FatalLogln("Error printing environment variables. Failed to convert the response")
	}

	if len(envs) == 0 {
		logger.SuccessLogln("No environment variable has been set yet")
	} else {
		logger.SuccessLogln("Following environment variables are set:")
		logger.Logln("")
		for k, v := range envs {
			s := fmt.Sprintf("%s=%s\n", k, v.(string))
			logger.Log(s)
		}
	}
}
