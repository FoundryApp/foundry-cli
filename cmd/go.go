package cmd

// "foundry go" or "foundry connect" or "foundry " or "foundry start" or "foundry link"?

import (
  "bytes"
  "crypto/md5"
  "encoding/hex"
  "encoding/json"
  // "mime/multipart"
  // "path/filepath"

  "log"
  "net/http"
  "io/ioutil"
  "fmt"
  "time"
  "github.com/spf13/cobra"

  "foundry/cli/logger"
  "foundry/cli/rwatch"
  "foundry/cli/auth"
  "foundry/cli/zip"

  "github.com/gorilla/websocket"
)

var (
  lastArchiveChecksum = ""
  goCmd = &cobra.Command{
    Use:    "go",
    Short:  "Connect Foundry to your cloud environment and GO!",
    Long:   "",
    Run:    runGo,
  }
  start = time.Now()

  uploadStart time.Time
)

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
  // TODO: Client shouldn't be using Auth ID token
  // baseURL := "127.0.0.1:3500"
  baseURL := "ide.foundryapp.co"

  // wsScheme := "ws"
  wsScheme := "wss"
  wsURL := fmt.Sprintf("%s://%s/ws/%s", wsScheme, baseURL, token)

  // pingScheme := "http"
  pingScheme := "https"
  pingURL := fmt.Sprintf("%s://%s/ping", pingScheme, baseURL)

  // url := fmt.Sprintf("ws://127.0.0.1:3500/ws/%s", token)
  // url := fmt.Sprintf("wss://ide.foundryapp.co/ws/%s", token)

  // url = "ws://ide.foundryapp.co/ws/token"

  logger.Debugln("=WS DIALING=")

  c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("WS dial error:", err)
  }
  defer c.Close()

  logger.Debugln("=WS CONNECTED=")

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
        // upload(c, token)
        uploadBuffer(c)
      case err := <-w.Errors:
				log.Println("watcher error:", err)
      }
    }
  }()

  err = w.AddRecursive(conf.RootDir)
  if err != nil {
    log.Println(err)
  }

  // Start periodically pinging server so the env isn't killed
  ticker := time.NewTicker(time.Second * 10)
  go ping(ticker, token, pingURL)

  // Don't wait for first save to send the code - send it as soon
  // as user calls 'foundry go'
  // upload(c, token)
  uploadBuffer(c)

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

func ping(ticker *time.Ticker, token, url string) {
  for {
    select {
    case <- ticker.C:
      // Ping the server
      var body = struct {
        Token string `json:"token"`
      }{token}

      jBody, err := json.Marshal(body)
      if err != nil {
        logger.Debugln("Error marshaling ping body: ", err)
        continue
      }

      res, err := http.Post(url, "application/json", bytes.NewBuffer(jBody))
      if err != nil {
        logger.Debugln("Error making ping post request: ", err)
        continue
      }

      if res.StatusCode != http.StatusOK {
        bodyBytes, err := ioutil.ReadAll(res.Body)
        if err != nil {
          logger.Debugln("Error reading ping response body: ", err)
          continue
        }

        bodyString := string(bodyBytes)
        logger.Debugln("Non-OK ping response: %s\n", bodyString)
      }
    }
  }
}

func uploadBuffer(c *websocket.Conn) {
  logger.Debugf("\n[timer] Starting timer\n");
  uploadStart = time.Now()

  ignore := []string{"node_modules", ".git"}

  // Zip the project
  buf, err := zip.ArchiveDir(conf.RootDir, ignore)
  if err != nil {
    log.Fatal(err)
  }


  // Get the 16 bytes hash
  archiveChecksum := checksum(buf.Bytes())

  // TODO: REMOVE
  // if lastArchiveChecksum == archiveChecksum { return }
  lastArchiveChecksum = archiveChecksum

  bufferSize := 1024 // 1024B, size of a single chunk
  buffer := make([]byte, bufferSize)
  chunkCount := (buf.Len() / bufferSize) + 1

  checksum := [md5.Size]byte{}
  previousChecksum := [md5.Size]byte{}

  for i := 0; i < chunkCount; i++ {
    bytesread, err := buf.Read(buffer)
    if err != nil {
      log.Fatal(err)
    }

    previousChecksum = checksum
    bytes := buffer[:bytesread]
    checksum = md5.Sum(bytes)

    checkStr := hex.EncodeToString(checksum[:])
    prevCheckStr := hex.EncodeToString(previousChecksum[:])

    lastChunk := i == chunkCount - 1

    if (lastChunk) {
      elapsed := time.Since(uploadStart)
      logger.Debugf("[timer] time until last chunk - %v\n", elapsed);
    }

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
    Data              string   `json:"data"`
    PreviousChecksum  string   `json:"previousChecksum"`
    Checksum          string   `json:"checksum"`
    IsLast            bool     `json:"isLast"`
    RunAll            bool     `json:"runAll"`
    Run               []string `json:"run"`
  }{hex.EncodeToString(b),
    prevChecksum,
    checksum,
    last,
    true,
    []string{},
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

    elapsed := time.Since(uploadStart)
    logger.Debugf("[timer] time until response - %v\n\n", elapsed);

    if err != nil {
      elapsed := time.Since(start)
      logger.Debugf("Elapsed time %s\n", elapsed)

      log.Fatal("WS error:", err)
    }
    fmt.Printf("%s\n", msg)
  }
}

func checksum(data []byte) string {
  hashInBytes := md5.Sum(data)
  return hex.EncodeToString(hashInBytes[:])
}
