package cmd

// "foundry go" or "foundry connect" or "foundry " or "foundry start" or "foundry link"?

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"foundry/cli/auth"
	conn "foundry/cli/connection"
	connMsg "foundry/cli/connection/msg"
	"foundry/cli/files"
	"foundry/cli/logger"
	p "foundry/cli/prompt"
	promptCmd "foundry/cli/prompt/cmd"
	"foundry/cli/rwatch"

	"github.com/spf13/cobra"
)

var (
	lastArchiveChecksum = ""
	goCmd               = &cobra.Command{
		Use:   "go",
		Short: "Connect Foundry to your cloud environment and GO!",
		Long:  "",
		Run:   runGo,
	}

	prompt *p.Prompt
	df     *os.File
)

func init() {
	rootCmd.AddCommand(goCmd)
}

func runGo(cmd *cobra.Command, args []string) {
	logger.Log("\n")
	warningText := "You aren't signed in. Some features aren't available! To sign in, run \x1b[1m'foundry sign-in'\x1b[0m or \x1b[1m'foundry sign-up'\x1b[0m to sign up.\n"

	switch authClient.AuthState {
	case auth.AuthStateTypeSignedOut:
		// Sign in anonmoysly + notify user
		if err := authClient.SignUpAnonymously(); err != nil {
			logger.FdebuglnFatal(err)
			logger.ErrorLoglnFatal(err)
		}

		if authClient.Error != nil {
			logger.FdebuglnFatal(authClient.Error)
			logger.ErrorLoglnFatal(authClient.Error)
		}

		logger.WarningLogln(warningText)
		time.Sleep(time.Second)
	case auth.AuthStateTypeSignedInAnonymous:
		// Notify user
		logger.WarningLogln(warningText)
		time.Sleep(time.Second)
	}

	done := make(chan struct{})

	// Create a new connection to the cloud env
	c, err := conn.New(authClient.IDToken)
	if err != nil {
		logger.FdebuglnFatal("Connection error", err)
		logger.ErrorLoglnFatal(err)
	}
	defer c.Close()

	watchCmd := promptCmd.NewWatchCmd()
	exitCmd := promptCmd.NewExitCmd()
	cmds := []promptCmd.Cmd{watchCmd, exitCmd}
	prompt = p.NewPrompt(cmds)
	go prompt.Run()

	// Listen for messages from the WS connection
	go c.Listen(listenCallback)

	// Start periodically pinging server so the env isn't killed
	pingMsg := connMsg.NewPingMsg(conn.PingURL(), authClient.IDToken)
	ticker := time.NewTicker(time.Second * 10)
	go c.Ping(pingMsg, ticker, done)

	// Start the file watcher
	w, err := rwatch.New(foundryConf.Ignore)
	if err != nil {
		logger.FdebuglnFatal("Watcher error", err)
		logger.ErrorLoglnFatal(err)
	}
	defer w.Close()

	err = w.AddRecursive(foundryConf.RootDir)
	if err != nil {
		logger.FdebuglnFatal("watcher AddRecursive", err)
		logger.ErrorLoglnFatal(err)
	}

	initialUploadCh := make(chan struct{}, 1)

	// The main goroutine handling all file events + prompt command requests
	// Command requests are all handled from a single goroutine because
	// Gorilla's websocket connection supports only one concurrent reader
	// and one concurrent writer - https://godoc.org/github.com/gorilla/websocket#hdr-Concurrency
	go func() {
		for {
			select {
			case args := <-watchCmd.RunCh:
				watchCmd.Run(c, args)
			case args := <-exitCmd.RunCh:
				exitCmd.Run(c, args)
			case <-initialUploadCh:
				files.Upload(c, foundryConf.RootDir, foundryConf.Ignore...)
			case _ = <-w.Events:
				// log.Println(e)
				files.Upload(c, foundryConf.RootDir, foundryConf.Ignore...)
			case err := <-w.Errors:
				logger.FdebuglnFatal("File watcher error", err)
				logger.ErrorLoglnFatal("File watcher error", err)
			}
		}
	}()

	// Don't wait for the first save event to send the code.
	// Send it as soon as user calls 'foundry go'
	initialUploadCh <- struct{}{}

	<-done
}

func listenCallback(data []byte, err error) {
	logger.Fdebugln(string(data))

	if err != nil {
		logger.FdebuglnFatal("WebSocket error", err)
		logger.ErrorLoglnFatal("WebSocket error", err)
	}

	t := connMsg.ResponseMsgType{}
	if err := json.Unmarshal(data, &t); err != nil {
		logger.FdebuglnFatal("Unmarshaling response error", err)
		logger.ErrorLoglnFatal("Parsing server JSON response error", err)
	}

	switch t.Type {
	case connMsg.LogResponseMsg:
		var s struct{ Content connMsg.LogContent }

		if err := json.Unmarshal(data, &s); err != nil {
			logger.FdebuglnFatal("Unmarshaling response error", err)
			logger.ErrorLoglnFatal("Parsing server log message error", err)
		}

		if _, err := prompt.Writeln(s.Content.Msg); err != nil {
			logger.FdebuglnFatal("Error writing output", err)
			logger.ErrorLoglnFatal("Error writing output", err)
		}

		// s1 := fmt.Sprintf("[0] %s", s.Content.Msg)
		// s2 := fmt.Sprintf("[1] %s", s.Content.Msg)
		// s3 := fmt.Sprintf("[2] %s\n", s.Content.Msg)
		// s4 := fmt.Sprintf("[3] %s\n", s.Content.Msg)
		// s5 := fmt.Sprintf("[4] %s\n", s.Content.Msg)
		// if _, err := prompt.Writeln(s1); err != nil {
		// 	logger.FdebuglnFatal(err)
		// }

	case connMsg.WatchResponseMsg:
		var s struct{ Content connMsg.WatchContent }

		if err := json.Unmarshal(data, &s); err != nil {
			logger.FdebuglnFatal("Unmarshaling response error", err)
			logger.ErrorLoglnFatal("Parsing server wathc message error", err)
		}

		p := fmt.Sprintf("Watching only functions: %s", strings.Join(s.Content.Run, ", "))
		prompt.SetInfoln(p)
	}
}
