package cmd

// "foundry go" or "foundry connect" or "foundry " or "foundry start" or "foundry link"?

import (
  "time"
  "fmt"

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
)

func init() {
  rootCmd.AddCommand(goCmd)
}

func runGo(cmd *cobra.Command, args []string) {
  token, err := getToken()
  if err != nil {
    logger.LogFatal("getToken error", err)
  }

  done := make(chan struct{})

  // Create a new connection to the cloud env
  c, err := conn.New(token)
  if err != nil {
    logger.LogFatal("Connection error", err)
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
    logger.LogFatal("Watcher error", err)
  }
  defer w.Close()

  go func() {
    for {
      select {
      case _ = <-w.Events:
        // log.Println(e)
        logger.Debugln("<timer> reseting starting upload time");
        uploadStart = time.Now()
        files.Upload(c, conf.RootDir)
      case err := <-w.Errors:
				logger.Debugln("watcher error:", err)
      }
    }
  }()

  err = w.AddRecursive(conf.RootDir)
  if err != nil {
    logger.Debugln(err)
  }

  // Don't wait for first save to send the code - send it as soon
  // as user calls 'foundry go'
  logger.Debugln("<timer> reseting starting upload time");
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
  logger.Debugf("<timer> time until response - %v\n", elapsed);

  if err != nil {
    elapsed := time.Since(start)
    logger.Debugf("Elapsed time %s\n", elapsed)

    logger.LogFatal("WS error:", err)
  }

  fmt.Printf("%s\n", data)
}
