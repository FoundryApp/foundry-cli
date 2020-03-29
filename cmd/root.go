package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"foundry/cli/auth"
	"foundry/cli/logger"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type FoundryConf struct {
	RootDir string `yaml:"rootDir"`
}

const confFile = "./foundry.config.yaml"

var (
	debugFile = ""
	authClient = auth.New()

	conf = FoundryConf{}
	rootCmd = &cobra.Command{
		Use:   "foundry",
		Short: "Better serverless dev",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			log.Println("Root command - add my description and implementation!")
		},
	}
)

func cobraInitCallback() {
	logger.InitDebug(debugFile)

	authClient.LoadTokens()
	if err := authClient.RefreshIDToken(); err != nil {
    logger.FdebuglnFatal(err)
  }
}

func init() {
	cobra.OnInitialize(cobraInitCallback)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&debugFile, "debug-file", "", "A file where the debug logs are saved (required)")

	confData, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Fatal("Read file error", err)
	}

	err = yaml.Unmarshal(confData, &conf)
	if err != nil {
		log.Fatal("YAML error", err)
	}
}

func Execute() {
	logger.Close()
	if err := rootCmd.Execute(); err != nil {
    log.Println(err)
    os.Exit(1)
	}
}