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
	res, err := firebase.Call("getUserEnvs", authClient.IDToken, nil)
	if err != nil {
		logger.FdebuglnFatal("Error calling getUserEnvs:", err)
		logger.FatalLogln("Error printing environment variables (1):", err)
	}
	if res.Error != nil {
		logger.FdebuglnFatal("Error calling getUserEnvs:", res.Error)
		logger.FatalLogln("Error printing environment variables (2):", res.Error)
	}

	envs, ok := res.Result.(map[string]interface{})
	if !ok {
		logger.FdebuglnFatal("Failed to type assert res.Result")
		logger.FatalLogln("Error printing environment variables. Failed to convert the resonse")
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
