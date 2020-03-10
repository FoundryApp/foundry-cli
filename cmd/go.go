package cmd

// "foundry go" or "foundry connect" or "foundry " or "foundry start" or "foundry link"?

import (
  "bytes"
  "mime/multipart"
  "path/filepath"

  "io"
  "log"
  "os"
  "net/http"
  // "io/ioutil"
  "fmt"
  "github.com/spf13/cobra"
  // "github.com/fsnotify/fsnotify"

  "foundry/cli/rwatch"
  "foundry/cli/auth"
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
  token, err := getToken()
  if err != nil {
    log.Fatal("getToken error", err)
  }

  w, err := rwatch.New()
  if err != nil {
    log.Fatal("Watcher error", err)
  }
  defer w.Close()

  done := make(chan bool)
  go func() {
    for {
      select {
      case _ = <-w.Events:
        // log.Println(e)
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

func getToken() (string, error) {
  a := auth.New()
  a.LoadTokens()
  if err := a.RefreshIDToken(); err != nil {
    return "", err
  }
  return a.IDToken, nil
}

func upload(token string) {
  ignore := []string{"node_modules", ".git", ".foundry"}
  // Zip project
  path, err := zip.ArchiveDir(conf.RootDir, ignore)
  if err != nil {
    log.Println("error", err)
  }

  // Send to cloud

  req, err := newFileUploadReq(path, token)
   if err != nil {
    log.Println("error file upload req", err)
  }
  client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
    if err != nil {
			log.Fatal(err)
		}
    resp.Body.Close()
		// fmt.Println(resp.StatusCode)
		// fmt.Println(resp.Header)
		fmt.Println(body)
	}
}

func newFileUploadReq(path, token string) (*http.Request, error) {
  file, err := os.Open(path)
  if err != nil {
		return nil, err
	}
  defer file.Close()

  body := &bytes.Buffer{}
  writer := multipart.NewWriter(body)
  part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, err
  }

  _, err = io.Copy(part, file)
  writer.WriteField("token", token)
  err = writer.Close()
  if err != nil {
		return nil, err
  }

  url := "http://127.0.0.1:8081/run"
  // url := "https://ide.foundryapp.co/run"

  req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}