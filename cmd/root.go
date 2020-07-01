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

// func LoadYAMLConfig() {
// 	fmt.Println("Loading foundry.yaml...")

// 	if _, err := os.Stat(confFile); os.IsNotExist(err) {
// 		logger.DebuglnError("Foundry config file 'foundry.yaml' not found in the current directory")
// 		logger.FatalLogln("Foundry config file 'foundry.yaml' not found in the current directory")
// 	}

// 	confData, err := ioutil.ReadFile(confFile)
// 	if err != nil {
// 		logger.DebuglnError("Can't read 'foundry.yaml' file", err)
// 		logger.FatalLogln("Can't read 'foundry.yaml' file", err)
// 	}

// 	err = yaml.Unmarshal(confData, &foundryConf)
// 	if err != nil {
// 		logger.DebuglnError("Config file 'foundry.yaml' isn't valid", err)
// 		logger.FatalLogln("Config file 'foundry.yaml' isn't valid", err)
// 	}

// 	dir, err := os.Getwd()
// 	if err != nil {
// 		logger.DebuglnError("Couldn't get current working directory", err)
// 		logger.FatalLogln("Couldn't get current working directory", err)
// 	}
// 	foundryConf.CurrentDir = dir

// 	// Parse IgnoreStr to globs
// 	for _, p := range foundryConf.IgnoreStrPatterns {
// 		// Add foundryConf.CurrentDir as a prefix to every glob pattern so
// 		// the prefix is same with file paths from watcher and  zipper

// 		// last := foundryConf.RootDir[len(foundryConf.RootDir)-1:]
// 		// if last != string(os.PathSeparator) {
// 		// 	p = foundryConf.RootDir + string(os.PathSeparator) + p
// 		// } else {
// 		// 	p = foundryConf.RootDir + p
// 		// }

// 		p = filepath.Join(foundryConf.CurrentDir, p)
// 		g, err := glob.Compile(p)
// 		if err != nil {
// 			logger.DebuglnError("Invalid glob pattern in the 'ignore' field in the foundry.yaml file")
// 			logger.FatalLogln("Invalid glob pattern in the 'ignore' field in the foundry.yaml file")
// 		}
// 		foundryConf.Ignore = append(foundryConf.Ignore, g)
// 	}
// }

// func CreateYAMLConfig() {
// 	logger.Debugln("Will try to create new foundry.yaml")

// 	if _, err := os.Stat(confFile); !os.IsNotExist(err) {
// 		logger.Debugln("foundry.yaml already exists")
// 		return
// 	}

// 	dest := filepath.Join(foundryConf.CurrentDir, "foundry.yaml")
// 	err := ioutil.WriteFile(dest, []byte(getYAMLConfigContent()), 0644)
// 	if err != nil {
// 		logger.FdebuglnError("Error writing foundry.yaml:", err)
// 		logger.FatalLogln("Error creating Foundry config file 'foundry.yaml':", err)
// 	}
// }

// func getYAMLConfigContent() string {
// 	return `
// # An array of glob patterns for files that should be ignored. The path is relative to the root dir.
// # If the array is changed, the CLI must be restarted for it to take the effect
// # See https://docs.foundryapp.co/configuration-file/ignore-directories-or-files
// ignore:
//     # Skip the whole node_modules directory
//   - node_modules
//     # Skip the whole .git directory
//   - .git
//     # Skip all hidden files
//   - "**/.*"
//     # Skip vim's temp files
//   - "**/*~"
//     # Ignore Firebase log files
//   - "**/firebase-*.log"

// # To access production data specify a path of a service account to your Firebase project
// # See https://docs.foundryapp.co/configuration-file/config-production-data
// # serviceAcc: path/to/your/service/acc/json/file
// `
// }
