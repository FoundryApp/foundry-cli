package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"foundry/cli/auth"
	conn "foundry/cli/connection"
	"foundry/cli/logger"

	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type FoundryConf struct {
	ServiceAccPath    string   `yaml:"serviceAcc"`
	IgnoreStrPatterns []string `yaml:"ignore"`

	CurrentDir string      `yaml:"-"` // Current working directory of CLI
	Ignore     []glob.Glob `yaml:"-"`
}

// Search a Foundry config file in the same directory from what was the foundry CLI called
const confFile = "./foundry.yaml"

var (
	debugFile        = ""
	authClient       *auth.Auth
	connectionClient *conn.Connection
	foundryConf      = FoundryConf{}

	rootCmd = &cobra.Command{
		Use:     "foundry",
		Short:   "Better serverless dev",
		Example: "foundry --help",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Logln("Foundry v0.2.3\n")
			logger.Logln("No subcommand was specified. To see all commands type 'foundry --help'")
		},
	}
)

func init() {
	// WARNING: logger's debug file isn't initialized yet. We can log only to the stdout or stderr.

	if len(os.Args) == 1 {
		return
	}

	cmd := os.Args[1]

	cobra.OnInitialize(func() { cobraInitCallback(cmd) })

	AddRootFlags(rootCmd)

	// TODO: Can this be in cobraInitCallback instead of here?
	if cmd != "init" &&
		cmd != "sign-out" &&
		cmd != "sign-in" &&
		cmd != "sign-up" &&
		cmd != "env-set" &&
		cmd != "env-print" &&
		cmd != "--help" {
		fmt.Println("Loading foundry.yaml...")

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
			logger.DebuglnError("Config file 'foundry.yaml' isn't valid", err)
			logger.FatalLogln("Config file 'foundry.yaml' isn't valid", err)
		}

		dir, err := os.Getwd()
		if err != nil {
			logger.DebuglnError("Couldn't get current working directory", err)
			logger.FatalLogln("Couldn't get current working directory", err)
		}
		foundryConf.CurrentDir = dir

		// Parse IgnoreStr to globs
		for _, p := range foundryConf.IgnoreStrPatterns {
			// Add foundryConf.CurrentDir as a prefix to every glob pattern so
			// the prefix is same with file paths from watcher and  zipper

			// last := foundryConf.RootDir[len(foundryConf.RootDir)-1:]
			// if last != string(os.PathSeparator) {
			// 	p = foundryConf.RootDir + string(os.PathSeparator) + p
			// } else {
			// 	p = foundryConf.RootDir + p
			// }

			p = filepath.Join(foundryConf.CurrentDir, p)
			g, err := glob.Compile(p)
			if err != nil {
				logger.DebuglnError("Invalid glob pattern in the 'ignore' field in the foundry.yaml file")
				logger.FatalLogln("Invalid glob pattern in the 'ignore' field in the foundry.yaml file")
			}
			foundryConf.Ignore = append(foundryConf.Ignore, g)
		}
	}
}

func cobraInitCallback(cmd string) {
	if err := logger.InitDebug(debugFile); err != nil {
		logger.DebuglnFatal("Failed to initialize a debug file for logger", err)
	}

	a, err := auth.New()
	if err != nil {
		logger.FdebuglnError("Error initializing Auth", err)
		logger.FatalLogln("Error initializing Auth", err)
	}
	if err := a.RefreshIDToken(); err != nil {
		logger.FdebuglnError("Error refreshing ID token: ", err)
		logger.FatalLogln("Error refreshing ID token: ", err)
	}
	authClient = a

	if cmd != "init" &&
		cmd != "sign-out" &&
		cmd != "sign-in" &&
		cmd != "sign-up" &&
		cmd != "--help" {
		logger.Log("\n")
		warningText := "You aren't signed in. Some features won't be available! To sign in, run \x1b[1m'foundry sign-in'\x1b[0m or \x1b[1m'foundry sign-up'\x1b[0m to sign up.\nThis message will self-destruct in 5s...\n"

		// Check if user signed in
		switch authClient.AuthState {
		case auth.AuthStateTypeSignedOut:
			// Sign in anonmoysly + notify user
			if err := authClient.SignUpAnonymously(); err != nil {
				logger.FdebuglnFatal(err)
				logger.FatalLogln(err)
			}

			if authClient.Error != nil {
				logger.FdebuglnFatal(authClient.Error)
				logger.FatalLogln(authClient.Error)
			}

			logger.WarningLogln(warningText)
			time.Sleep(time.Second * 5)
		case auth.AuthStateTypeSignedInAnonymous:
			// Notify user
			logger.WarningLogln(warningText)
			time.Sleep(time.Second)
		}

		// TODO: Now only 'go' command can use connectionClient variable
		// This should be handled better
		if cmd == "go" {
			// Create a new connection to the cloud env
			fmt.Println("Connecting to your cloud environment...")
			c, err := conn.New(authClient.IDToken)
			if err != nil {
				logger.FdebuglnFatal("Connection error", err)
				logger.FatalLogln(err)
			}
			connectionClient = c
		}
	}
}

func Execute() {
	defer func() {
		if connectionClient != nil {
			connectionClient.Close()
		}
		logger.Close()
	}()

	if err := rootCmd.Execute(); err != nil {
		logger.FdebuglnError("Error executing root command", err)
		logger.FatalLogln(err)
	}
}
