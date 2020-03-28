package cmd

// "foundry go" or "foundry connect" or "foundry " or "foundry start" or "foundry link"?

import (
  "time"
  "fmt"
  "os"
  "encoding/json"
  "strings"

  "foundry/cli/auth"
  conn "foundry/cli/connection"
  connMsg "foundry/cli/connection/msg"
  pc "foundry/cli/prompt/cmd"
  "foundry/cli/files"
  "foundry/cli/logger"
  p "foundry/cli/prompt"
  "foundry/cli/rwatch"

  "github.com/spf13/cobra"
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

  prompt *p.Prompt
  uploadStart = time.Now()

  df *os.File
)

func init() {
  rootCmd.AddCommand(goCmd)
}

func runGo(cmd *cobra.Command, args []string) {
  token, err := getToken()
  if err != nil {
    logger.FdebuglnFatal("getToken error", err)
    logger.LogFatal(err)
  }

  done := make(chan struct{})

  // Create a new connection to the cloud env
  c, err := conn.New(token)
  if err != nil {
    logger.FdebuglnFatal("Connection error", err)
    logger.LogFatal(err)
  }
  defer c.Close()

  // Listen for messages from the WS connection
  go c.Listen(listenCallback)

  // Start periodically pinging server so the env isn't killed
  pingMsg := connMsg.NewPingMsg(conn.PingURL(), token)
  ticker := time.NewTicker(time.Second * 10)
  go c.Ping(pingMsg, ticker, done)

  // Start an interactive prompt
  cmds := []*p.Cmd{
    pc.Watch(c),
    pc.Exit(),
  }
  prompt = p.NewPrompt(cmds)
  go prompt.Run()

  // Start the file watcher
  w, err := rwatch.New()
  if err != nil {
    logger.FdebuglnFatal("Watcher error", err)
    logger.LogFatal(err)
  }
  defer w.Close()

  go func() {
    for {
      select {
      case _ = <-w.Events:
        // log.Println(e)
        logger.Fdebugln("<timer> reseting starting upload time")
        uploadStart = time.Now()
        files.Upload(c, conf.RootDir)
      case err := <-w.Errors:
        logger.FdebuglnFatal("watcher error", err)
        logger.LogFatal(err)
      }
    }
  }()

  err = w.AddRecursive(conf.RootDir)
  if err != nil {
    logger.FdebuglnFatal("watcher AddRecursive", err)
    logger.LogFatal(err)
  }

  // Don't wait for first save to send the code - send it as soon
  // as user calls 'foundry go'
  logger.Fdebugln("<timer> reseting starting upload time")
  uploadStart = time.Now()
  files.Upload(c, conf.RootDir)

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

func listenCallback(data []byte, err error) {
  elapsed := time.Since(uploadStart)
  logger.Fdebugln("<timer> time until response (+ time.Sleep) -", elapsed)

  // time.Sleep(time.Millisecond * 20)
  logger.Fdebugln(string(data))



  if err != nil {
    elapsed := time.Since(start)
    logger.Fdebugln("<timer> Elapsed time -", elapsed)
    logger.FdebuglnFatal("WS error", err)
    logger.LogFatal(err)
  }


  t := connMsg.ResponseMsgType{}
  if err := json.Unmarshal(data, &t); err != nil {
    elapsed := time.Since(start)
    logger.Fdebugln("<timer> Elapsed time -", elapsed)
    logger.FdebuglnFatal("Unmarshaling response error", err)
    logger.LogFatal(err)
  }

  switch t.Type {
  case connMsg.LogResponseMsg:
    var s struct { Content connMsg.LogContent }

    if err := json.Unmarshal(data, &s); err != nil {
      elapsed := time.Since(start)
      logger.Fdebugln("<timer> Elapsed time -", elapsed)
      logger.FdebuglnFatal("Unmarshaling response error", err)
    }

    // TODO: listenCallback is a callback - it doesn't wait for prompt to print everything
    // prompt must have a buffer and a lock that makes sure that it's printing sequentially
    // t := fmt.Sprintf("%s\n", s.Content.Msg)
    prompt.Print(string(s.Content.Msg))
    // if err := prompt.WriteToBuffer(t); err != nil {
    //   logger.FdebuglnFatal(err)
		// 	logger.LogFatal(err)
    // }

  case connMsg.WatchResponseMsg:
    var s struct { Content connMsg.WatchContent }

    if err := json.Unmarshal(data, &s); err != nil {
      elapsed := time.Since(start)
      logger.Fdebugln("<timer> Elapsed time -", elapsed)
      logger.FdebuglnFatal("Unmarshaling response error", err)
      logger.LogFatal(err)
    }

    p := fmt.Sprintf("[%s] > ", strings.Join(s.Content.Run, ", "))
    prompt.SetPromptPrefix(p)
  }
}
