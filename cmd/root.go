package cmd

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"foundry/cli/logger"

	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
)

type FoundryConf struct {
	CurrentDir string // current working directory of CLI
	Ignore     []glob.Glob
}

const packageJSONFile = "./package.json"
const ignoreFile = "./.foundryignore"

var (
	debugFile   = ""
	foundryConf = FoundryConf{}

	rootCmd = &cobra.Command{
		Use:     "foundry",
		Short:   "Better serverless dev",
		Example: "foundry --help",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Logln("No subcommand was specified. To see all commands type 'foundry --help'")
		},
	}
)

func init() {
	// WARNING: logger's debug file isn't initialized yet. We can log only to the stdout or stderr.
	cobra.OnInitialize(func() { cobraInitCallback() })
	AddRootFlags(rootCmd)
}

func cobraInitCallback() {
	if err := logger.InitDebug(debugFile); err != nil {
		logger.DebuglnFatal("Failed to initialize a debug file for logger", err)
	}
}

func Execute() {
	defer func() {
		logger.Close()
	}()

	dir, err := os.Getwd()
	if err != nil {
		logger.DebuglnError("Couldn't get current working directory", err)
		logger.FatalLogln("Couldn't get current working directory", err)
	}
	foundryConf.CurrentDir = dir

	if err := rootCmd.Execute(); err != nil {
		logger.FdebuglnError("Error executing root command", err)
		logger.FatalLogln(err)
	}
}

func CheckForPackageJSON() {
	if _, err := os.Stat(packageJSONFile); os.IsNotExist(err) {
		logger.DebuglnError("package.json not found in the current directory")
		logger.FatalLogln("package.json not found in the current directory. Start Foundry from the same directory where is your Cloud Function package.json file.")
	}
}

func LoadIgnoreFile() {
	if _, err := os.Stat(ignoreFile); os.IsNotExist(err) {
		// Ignore file doesn't exist. Nothing to load, we can exit.
		return
	} else if err != nil {
		logger.DebuglnError("Error while getting info about ignore file '.foundryignore'", err)
		logger.FatalLogln("Couldn't read ignore file '.foundryignore'")
	}

	f, err := os.Open(ignoreFile)
	if err != nil {
		logger.DebuglnError("Error while opening ignore file '.foundryignore'", err)
		logger.FatalLogln("Couldn't read ignore file '.foundryignore'", err)
	}

	// Read ignore file line by line
	reader := bufio.NewReader(f)
	line, err := readln(reader)
	for err == nil {
		if line != "" {
			relativePath := filepath.Join(foundryConf.CurrentDir, strings.TrimSpace(line))
			g, compileErr := glob.Compile(relativePath)
			if compileErr != nil {
				logger.DebuglnError("Invalid glob pattern in the '.foundryignore' file", compileErr)
				logger.FatalLogln("Invalid glob pattern in the '.foundryignore' file")
			}
			foundryConf.Ignore = append(foundryConf.Ignore, g)
		}
		line, err = readln(reader)
	}
}

func readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
