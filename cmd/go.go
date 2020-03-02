package cmd

// "foundry go" or "foundry connect" or "foundry " or "foundry start" or "foundry link"?

import (
  "log"
  "os"
  "net/http"
  "io/ioutil"
  "fmt"
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
        upload()
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

func upload() {
  ignore := []string{"node_modules", ".git", ".foundry"}
  // Zip project
  path, err := zip.ArchiveDir(conf.RootDir, ignore)
  if err != nil {
    log.Println("error", err)
  }

  // Send to cloud
  makeUploadReq(path)
}

func makeUploadReq(fname string) {
  log.Println("makeUploadReq", fname)
  file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	res, err := http.Post("http://127.0.0.1:8080/run", "binary/octet-stream", file)
	if err != nil {
		panic(err)
  }
	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
  fmt.Printf(string(message))
}
