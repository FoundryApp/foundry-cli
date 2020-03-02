package cmd

// "foundry go" or "foundry connect" or "foundry " or "foundry start" or "foundry link"?

import (
  "log"
  // "time"
  "github.com/spf13/cobra"
  // "github.com/fsnotify/fsnotify"

  "foundry/cli/rwatch"
  "foundry/cli/zip"
)

var goCmd = &cobra.Command{
  Use:   "go",
  Short: "Connect Foundry to your project and GO!",
  Long:  "",
  Run: run,
}

func init() {
  rootCmd.AddCommand(goCmd)
}

func run(cmd *cobra.Command, args []string) {
  w, err := rwatch.New()
  if err != nil {
    log.Fatal("Watcher error", err)
  }
  defer w.Close()

  done := make(chan bool)
  go func() {
    for {
      select {
      case e := <-w.Events:
        log.Println(e)
        sendFiles()
      case err := <-w.Errors:
				log.Println("error:", err)
      }
    }
  }()

  err = w.AddRecursive(conf.RootDir)
  if err != nil {
    log.Println(err)
  }
  <-done
}

func sendFiles() {
  ignore := []string{"node_modules", ".git", ".foundry"}
  // Zip project
  f, err := zip.ArchiveDir(conf.RootDir, ignore)
  if err != nil {
    log.Println("error", err)
  }
  // log.Println(f.Name())

  // Send to cloud
}