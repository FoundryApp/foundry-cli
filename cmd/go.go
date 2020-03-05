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
  Use:    "go",
  Short:  "Connect Foundry to your project and GO!",
  Long:   "",
  Run:    runGo,
}

func init() {
  rootCmd.AddCommand(goCmd)
}

func runGo(cmd *cobra.Command, args []string) {
  // Connect to pod
  token := getToken()

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
        upload(token)
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

func getToken() string {
  return ""
}

func upload(token string) {
  ignore := []string{"node_modules", ".git", ".foundry"}
  // Zip project
  path, err := zip.ArchiveDir(conf.RootDir, ignore)
  if err != nil {
    log.Println("error", err)
  }

  // Send to cloud
  makeUploadReq(path, token)
}

func makeUploadReq(fname string, token string) {
  log.Println("makeUploadReq", fname)
  file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

  url := "http://127.0.0.1:8080/run/"
  // url := fmt.Sprintf("http://127.0.0.1:8080/run/%v", token)
  // url := fmt.Sprintf("http://ide.foundryapp.co/run/%v", token)

  res, err := http.Post(url, "binary/octet-stream", file)
	if err != nil {
    log.Println(err)
    panic(err)
  }
  defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
  fmt.Printf(string(message))
}