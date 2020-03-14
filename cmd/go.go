package cmd

// "foundry go" or "foundry connect" or "foundry " or "foundry start" or "foundry link"?

import (
  // "bytes"
  "crypto/md5"
  "encoding/hex"
  // "mime/multipart"
  // "path/filepath"

  // "io"
  "log"
  "os"
  // "net/http"
  // "io/ioutil"
  // "fmt"
  "github.com/spf13/cobra"
  // "github.com/fsnotify/fsnotify"

  "foundry/cli/rwatch"
  "foundry/cli/auth"
  "foundry/cli/zip"

  "github.com/gorilla/websocket"
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

  // Connect to websocket
  c, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:3500/ws", nil)
	if err != nil {
		log.Fatal("WS dial error:", err)
	}
	defer c.Close()

  go listenWS(c)

  // Start file watcher
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

        // TODO: Send binary data with websocket
        upload(c, token)
      case err := <-w.Errors:
				log.Println("watcher error:", err)
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

func upload(c *websocket.Conn, token string) {
  ignore := []string{"node_modules", ".git", ".foundry"}

  // Zip the project
  path, err := zip.ArchiveDir(conf.RootDir, ignore)
  if err != nil {
    log.Fatal(err)
  }

  // Read file in chunks and send each chunk
  file, err := os.Open(path)
  if err != nil {
		log.Fatal(err)
	}
  defer file.Close()

  // fileInfo is needed to calculate how many chunks are in the file
  fileInfo, err := os.Stat(path)
  if err != nil {
		log.Fatal(err)
  }

  bufferSize := int64(1024) // 1024B, size of a single chunk
  buffer := make([]byte, bufferSize)
  chunkCount := (fileInfo.Size() / bufferSize) + 1

  checksum := [md5.Size]byte{}
  previousChecksum := [md5.Size]byte{}

  for i := int64(0); i < chunkCount; i++ {
    bytesread, err := file.Read(buffer)
    if err != nil {
      log.Fatal(err)
    }

    previousChecksum = checksum
    bytes := buffer[:bytesread]
    checksum = md5.Sum(bytes)

    checkStr := hex.EncodeToString(checksum[:])
    prevCheckStr := hex.EncodeToString(previousChecksum[:])

    lastChunk := i == chunkCount - 1

    log.Println("Size", bytesread)

    if err = sendChunk(
      c,
      bytes,
      checkStr,
      prevCheckStr,
      lastChunk); err != nil {
      log.Fatal(err)
    }
  }
}


func sendChunk(c *websocket.Conn, b []byte, checksum string, prevChecksum string, last bool) error {
  msg := struct {
    Data              string `json:"data"`
    PreviousChecksum  string `json:"previousChecksum"`
    Checksum          string `json:"checksum"`
    IsLast            bool   `json:"isLast"`
  }{hex.EncodeToString(b),
    prevChecksum,
    checksum,
    last,
  }
  err := c.WriteJSON(msg)
  if err != nil {
    return err
  }
  return nil
}

func listenWS(c *websocket.Conn) {
  for {
    _, msg, err := c.ReadMessage()
    if err != nil {
      log.Fatal("WS error:", err)
    }
    log.Printf("Autorun message: %s", msg)
  }
}
