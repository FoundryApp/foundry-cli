package cmd

import (
	// "fmt"
	// "foundry/cli/connection/msg"
	// "foundry/cli/firebase"
	// "foundry/cli/logger"
	// "os"
	// "strings"

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
	// if len(args) == 0 {
	// 	logger.WarningLogln("No envs to delete specified. Example usage: 'foundry env-delete ENV_1 ENV_2'")
	// 	os.Exit(0)
	// }

	// reqBody := struct {
	// 	Delete []string `json:"delete"`
	// }{args}

	// s := fmt.Sprintf("Deleting following env variables: '%s'...", strings.Join(args, ","))
	// logger.Logln(s)

	// res, err := firebase.Call("deleteUserEnvs", authClient.IDToken, reqBody)
	// if err != nil {
	// 	logger.FdebuglnFatal("Error calling deleteUserEnvs:", err)
	// 	logger.FatalLogln("Error deleting environment variables (1):", err)
	// }
	// if res.Error != nil {
	// 	logger.FdebuglnFatal("Error calling deleteUserEnvs:", res.Error)
	// 	logger.FatalLogln("Error deleting environment variables (2):", res.Error)
	// }

	// // Send new envs to Autorun
	// logger.Fdebugln("New env vars after deletion:", res.Result)

	// envsMap, ok := res.Result.(map[string]interface{})
	// if !ok {
	// 	logger.FdebuglnFatal("Failed to type assert res.Result")
	// 	logger.FatalLogln("Error deleting environment variables (3)")
	// }

	// envs := []msg.Env{}
	// for name, val := range envsMap {
	// 	envs = append(envs, msg.Env{name, val.(string)})
	// }
	// logger.Fdebugln("Sending new envs vars to Autorun:", envs)

	// envMsg := msg.NewEnvMsg(authClient.IDToken, envs)
	// if err := envMsg.Send(); err != nil {
	// 	logger.FdebuglnError("Failed to report new env vars (after deletion) to Autorun", err)
	// 	logger.DebuglnError("Error deleting environment variables (4)", err)
	// 	return
	// }

	// // Print new envs
	// logger.SuccessLogln("Deleted")
	// logger.Logln("---------------")

	// logger.Logln("")
	// logger.Logln("Env variables now:")
	// if len(envsMap) == 0 {
	// 	logger.Logln("There are no env variables set in your environment")
	// } else {
	// 	for k, v := range envsMap {
	// 		s := fmt.Sprintf("\t%s=%s\n", k, v.(string))
	// 		logger.Log(s)
	// 	}
	// }
}
