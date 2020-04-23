package cmd

import (
	"foundry/cli/logger"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:     "init",
		Short:   "Create the initial foundry.yaml config file",
		Example: "foundry init",
		Run:     runInit,
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(confFile); !os.IsNotExist(err) {
		logger.FdebuglnError("Foundry config file 'foundry.yaml' already exists")
		logger.FatalLogln("Foundry config file 'foundry.yaml' already exists")
	}

	dest := filepath.Join(foundryConf.CurrentDir, "foundry.yaml")
	err := ioutil.WriteFile(dest, []byte(getInitYaml()), 0644)
	if err != nil {
		logger.FdebuglnError("Error writing foundry.yaml:", err)
		logger.FatalLogln("Error creating Foundry config file 'foundry.yaml':", err)
	}

	logger.SuccessLogln("Config file 'foundry.yaml' created")
}

func getInitYaml() string {
	// TODO: Update to a final version of the init config yaml
	return `
# An array of glob patterns for files that should be ignored. The path is relative to the root dir.
# If the array is changed, the CLI must be restarted for it to take the effect
# See https://docs.foundryapp.co/configuration-file/ignore-directories-or-files
ignore:
    # Skip the whole node_modules directory
  - node_modules
    # Skip the whole .git directory
  - .git
    # Skip all hidden files
  - "**/.*"
    # Skip vim's temp files
  - "**/*~"
    # Ignore Firebase log files
	- "**/firebase-debug.log"

# Enable TypeScript
# See https://docs.foundryapp.co/resources/supported-languages#using-foundry-with-cloud-functions-in-typescript
# typescript: true

# An array describing emulated Firebase Auth users in your cloud environment
# See https://docs.foundryapp.co/configuration-file/emulate-users
#users:    
#  - id: user-id-1
#    # The 'data' field takes a JSON string
#   data: '{"email": "user-id-1-email@email.com"}'  


# An array describing emulated Firestore in your cloud environment
# See https://docs.foundryapp.co/configuration-file/emulate-firestore
#firestore:    
#  # You can describe your emulated Firestore either directly
#  - collection: workspaces
#    docs:
#      - id: ws-id-1
#        data: '{"userId": "user-id-1"}'


# An array describing your Firebase functions that should be evaluated by Foundry. 
# All described functions must be exported in the function's root index.js file.
# In this array, you describe how Foundry should trigger each function in every run.
# See https://docs.foundryapp.co/configuration-file/config-functions
#functions:
#  - name: myHttpsFunction
#    type: https
#    payload: '{"field":"value"}'
`
}
