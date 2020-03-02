package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type FoundryConf struct {
	RootDir string `yaml:"rootDir"`
}

const confFile = "./foundry.config.yaml"

var (
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

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Fatal("Read file error", err)
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal("YAML error", err)
	}
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    log.Println(err)
    os.Exit(1)
  }
}