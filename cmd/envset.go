package cmd

import (
	"strings"

	"foundry/cli/connection/msg"
	"foundry/cli/logger"

	"github.com/spf13/cobra"
)

var (
	envSetCmd = &cobra.Command{
		Use:     "env-set",
		Short:   "Set environment variable(s) in your cloud environment",
		Example: "foundry env-set env-name=env-value",
		Args:    cobra.MinimumNArgs(1),
		Run:     runEnvSet,
	}
)

func init() {
	rootCmd.AddCommand(envSetCmd)
}

func runEnvSet(cmd *cobra.Command, args []string) {
	envs := []msg.Env{}
	for _, env := range args {
		arr := strings.Split(env, "=")

		if len(arr) != 2 {
			logger.FdebuglnFatal("Error parsing environment variable:", env)
			logger.FatalLogln("Error parsing environment variable. Expected format 'env=value'. Got:", env)
		}

		name := arr[0]
		val := arr[1]

		if name == "" {
			logger.FdebuglnFatal("Error parsing environment variable - name is empty:", env)
			logger.FatalLogln("Error parsing environment variable. Expected format 'env=value'. Got;", env)
		}
		if val == "" {
			logger.FdebuglnFatal("Error parsing environment variable - val is empty:", env)
			logger.FatalLogln("Error parsing environment variable. Expected format 'env=value'. Got:", env)
		}

		envs = append(envs, msg.Env{name, val})
	}

	envMsg := msg.NewEnvMsg(authClient.IDToken, envs)
	if err := envMsg.Send(); err != nil {
		logger.FdebuglnError("Error setting environment variables:", err)
		logger.DebuglnError("Error setting environment variables:", err)
		return
	}
	logger.SuccessLogln("Variables Set")
}
