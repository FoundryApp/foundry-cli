package cmd

import (
	"io/ioutil"
	"os"

	"foundry/cli/auth"
	"foundry/cli/logger"

	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type FoundryConf struct {
	RootDir           string   `yaml:"rootDir"`
	IgnoreStrPatterns []string `yaml:"ignore"`

	Ignore []glob.Glob `yaml:"-"`
}

// Search a Foundry config file in the same directory from what was the foundry CLI called
const confFile = "./foundry.yaml"

var (
	debugFile  = ""
	authClient *auth.Auth

	foundryConf = FoundryConf{}
	rootCmd     = &cobra.Command{
		Use:   "foundry",
		Short: "Better serverless dev",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			logger.Logln("Root command - add my description and implementation!")
		},
	}
)

func cobraInitCallback() {
	if err := logger.InitDebug(debugFile); err != nil {
		logger.DebuglnFatal("Failed to initialized debug file for logger")
	}

	a, err := auth.New()
	if err != nil {
		logger.FdebuglnError("Error initializing Auth", err)
		logger.FatalLogln("Error initializing Auth", err)
	}
	if err := a.RefreshIDToken(); err != nil {
		logger.FdebuglnError("Error refreshing ID token", err)
		logger.FatalLogln("Error refreshing ID token", err)
	}
	authClient = a
}

func init() {
	// WARNING: logger's debug file isn't initialized yet. We can log only to the stdout or stderr.

	cobra.OnInitialize(cobraInitCallback)

	// DEBUG:
	rootCmd.PersistentFlags().StringVar(&debugFile, "debug-file", "", "A file where the debug logs are saved (required)")

	cmd := os.Args[1]
	if cmd != "init" {
		if _, err := os.Stat(confFile); os.IsNotExist(err) {
			logger.DebuglnError("Foundry config file 'foundry.yaml' not found in the current directory")
			logger.FatalLogln("Foundry config file 'foundry.yaml' not found in the current directory. Run '\x1b[1mfoundry init\x1b[0m'.")
		}

		confData, err := ioutil.ReadFile(confFile)
		if err != nil {
			logger.DebuglnError("Can't read 'foundry.yaml' file", err)
			logger.FatalLogln("Can't read 'foundry.yaml' file", err)
		}

		err = yaml.Unmarshal(confData, &foundryConf)
		if err != nil {
			logger.DebuglnError("foundry.yaml file isn't a valid YAML file or doesn't contain field 'RootDir'", err)
			logger.FatalLogln("foundry.yaml file isn't a valid YAML file or doesn't contain field 'RootDir'", err)
		}
		if foundryConf.RootDir == "" {
			logger.DebuglnError("foundry.yaml doesn't contain field 'RootDir' or it's empty")
			logger.FatalLogln("foundry.yaml doesn't contain field 'RootDir' or it's empty")
		}

		// Parse IgnoreStr to globs
		for _, p := range foundryConf.IgnoreStrPatterns {
			// Add "./" as a prefix to every glob pattern so
			// the prefix is same with file paths from watcher
			// zipper
			// last := foundryConf.RootDir[len(foundryConf.RootDir)-1:]
			// if last != string(os.PathSeparator) {
			// 	p = foundryConf.RootDir + string(os.PathSeparator) + p
			// } else {
			// 	p = foundryConf.RootDir + p
			// }

			g, err := glob.Compile(p)
			if err != nil {
				logger.DebuglnError("Invalid glob pattern in the 'ignore' field in the foundry.yaml file")
				logger.FatalLogln("Invalid glob pattern in the 'ignore' field in the foundry.yaml file")
			}
			foundryConf.Ignore = append(foundryConf.Ignore, g)
		}

		logger.Debugln("Ignore str", foundryConf.IgnoreStrPatterns)
		logger.Debugln("Ignore glob", foundryConf.Ignore)
	}

}

func Execute() {
	defer logger.Close()
	if err := rootCmd.Execute(); err != nil {
		logger.FdebuglnError("Error executing root command", err)
		logger.FatalLogln("Error executing root command", err)
	}
}
