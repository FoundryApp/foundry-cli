package cmd

import (
	"io/ioutil"
	"os"

	"foundry/cli/auth"
	"foundry/cli/logger"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type FoundryConf struct {
	RootDir string `yaml:"rootDir"`
}

const confFile = "./foundry.yaml"

var (
	debugFile  = ""
	authClient *auth.Auth

	conf    = FoundryConf{}
	rootCmd = &cobra.Command{
		Use:   "foundry",
		Short: "Better serverless dev",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			logger.Logln("Root command - add my description and implementation!")
		},
	}
)

func cobraInitCallback() {
	logger.InitDebug(debugFile)

	a, err := auth.New()
	if err != nil {
		logger.FdebuglnFatal(err)
	}
	if err := a.RefreshIDToken(); err != nil {
		logger.FdebuglnFatal(err)
	}
	authClient = a
}

func init() {
	cobra.OnInitialize(cobraInitCallback)

	// DEBUG:
	rootCmd.PersistentFlags().StringVar(&debugFile, "debug-file", "", "A file where the debug logs are saved (required)")

	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		logger.Fdebugln("Foundry config file 'foundry.yaml' not found in the current directory")
		logger.ErrorLoglnFatal("Foundry config file 'foundry.yaml' not found in the current directory")
	}

	confData, err := ioutil.ReadFile(confFile)
	if err != nil {
		logger.Fdebugln("Can't read 'foundry.yaml' file", err)
		logger.ErrorLoglnFatal("Can't read 'foundry.yaml' file", err)
	}

	err = yaml.Unmarshal(confData, &conf)
	if err != nil {
		logger.Fdebugln("foundry.yaml file isn't a valid YAML file or doesn't contain field 'RootDir'", err)
		logger.ErrorLoglnFatal("foundry.yaml file isn't a valid YAML file or doesn't contain field 'RootDir'", err)
	}
	if conf.RootDir == "" {
		logger.Fdebugln("foundry.yaml doesn't contain field 'RootDir' or it's empty")
		logger.ErrorLoglnFatal("foundry.yaml doesn't contain field 'RootDir' or it's empty")
	}
}

func Execute() {
	logger.Close()
	if err := rootCmd.Execute(); err != nil {
		logger.Fdebugln(err)
		logger.ErrorLoglnFatal(err)
	}
}
