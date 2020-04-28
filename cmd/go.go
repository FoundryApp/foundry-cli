package cmd

// "foundry go" or "foundry connect" or "foundry " or "foundry start" or "foundry link"?

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	conn "foundry/cli/connection"
	connMsg "foundry/cli/connection/msg"
	"foundry/cli/files"
	"foundry/cli/logger"
	p "foundry/cli/prompt"
	promptCmd "foundry/cli/prompt/cmd"
	"foundry/cli/rwatch"

	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
)

var (
	lastArchiveChecksum = ""
	goCmd               = &cobra.Command{
		Use:     "go",
		Short:   "Connect to your cloud environment and start watching your Firebase Functions",
		Example: "foundy go",
		Run:     runGo,
	}

	prompt *p.Prompt
	df     *os.File
)

func init() {
	rootCmd.AddCommand(goCmd)
}

func runGo(cmd *cobra.Command, args []string) {
	done := make(chan struct{})

	watchCmd := promptCmd.NewWatchCmd()
	watchAllCmd := promptCmd.NewWatchAllCmd()
	exitCmd := promptCmd.NewExitCmd()
	envPrintCmd := promptCmd.NewEnvPrintCmd(authClient.IDToken)
	envSetCmd := promptCmd.NewEnvSetCmd(authClient.IDToken)
	envDelCmd := promptCmd.NewEnvDelCmd(authClient.IDToken)

	cmds := []promptCmd.Cmd{watchCmd, watchAllCmd, exitCmd, envPrintCmd, envSetCmd, envDelCmd}
	prompt = p.NewPrompt(cmds)
	go prompt.Run()

	// Listen for messages from the WS connection
	go connectionClient.Listen(listenCallback)

	// Start periodically pinging server so the env isn't killed
	pingMsg := connMsg.NewPingMsg(conn.PingURL(), authClient.IDToken)
	ticker := time.NewTicker(time.Second * 10)
	go connectionClient.Ping(pingMsg, ticker, done)

	// Start the file watcher
	w, err := rwatch.New(foundryConf.Ignore)
	if err != nil {
		logger.FdebuglnFatal("Watcher error", err)
		logger.FatalLogln(err)
	}
	defer w.Close()

	err = w.AddRecursive(foundryConf.CurrentDir)
	if err != nil {
		logger.FdebuglnFatal("watcher AddRecursive", err)
		logger.FatalLogln(err)
	}

	initialUploadCh := make(chan struct{}, 1)
	promptNotifCh := make(chan string)

	// The main goroutine handling all file events + prompt command requests
	// Command requests are all handled from a single goroutine because
	// Gorilla's websocket connection supports only one concurrent reader
	// and one concurrent writer.
	// More info - https://godoc.org/github.com/gorilla/websocket#hdr-Concurrency
	go func() {
		for {
			select {
			case event := <-prompt.Events:
				if event.Type == p.PromptEventTypeRerender {
					_ = prompt.ShowLoading()
					files.Upload(connectionClient, foundryConf.CurrentDir, foundryConf.ServiceAccPath, promptNotifCh, foundryConf.Ignore...)
				}
			case msg := <-promptNotifCh:
				prompt.SetInfoln(msg, p.InfoLineSeverityNormal)
			case args := <-envDelCmd.RunCh:
				_ = prompt.ShowLoading()
				pOut, pInfo, err := envDelCmd.Run(connectionClient, args)
				if err != nil {
					prompt.SetInfoln(err.Error(), p.InfoLineSeverityError)
					continue
				}
				prompt.SetInfoln(pInfo, p.InfoLineSeverityNormal)
				prompt.Writeln(pOut)
			case args := <-envSetCmd.RunCh:
				_ = prompt.ShowLoading()
				pOut, pInfo, err := envSetCmd.Run(connectionClient, args)
				if err != nil {
					prompt.SetInfoln(err.Error(), p.InfoLineSeverityError)
					continue
				}
				prompt.SetInfoln(pInfo, p.InfoLineSeverityNormal)
				prompt.Writeln(pOut)
			case args := <-envPrintCmd.RunCh:
				_ = prompt.ShowLoading()
				pOut, pInfo, err := envPrintCmd.Run(connectionClient, args)
				if err != nil {
					prompt.SetInfoln(err.Error(), p.InfoLineSeverityError)
					continue
				}
				prompt.SetInfoln(pInfo, p.InfoLineSeverityNormal)
				prompt.Writeln(pOut)
			case args := <-watchAllCmd.RunCh:
				if _, _, err := watchAllCmd.Run(connectionClient, args); err != nil {
					prompt.SetInfoln(err.Error(), p.InfoLineSeverityError)
					continue
				}
			case args := <-watchCmd.RunCh:
				_, pInfo, err := watchCmd.Run(connectionClient, args)
				if err != nil {
					prompt.SetInfoln(err.Error(), p.InfoLineSeverityError)
					continue
				}
				prompt.SetInfoln(pInfo, p.InfoLineSeverityError)
			case args := <-exitCmd.RunCh:
				_, _, _ = exitCmd.Run(connectionClient, args)
			case <-initialUploadCh:
				files.Upload(connectionClient, foundryConf.CurrentDir, foundryConf.ServiceAccPath, promptNotifCh, foundryConf.Ignore...)
			case e := <-w.Events:
				path := "." + string(os.PathSeparator) + e.Name
				if !ignored(path, foundryConf.Ignore) {
					logger.Fdebugln("Watcher event", e.Name)
					_ = prompt.ShowLoading()
					files.Upload(connectionClient, foundryConf.CurrentDir, foundryConf.ServiceAccPath, promptNotifCh, foundryConf.Ignore...)
				}
			case err := <-w.Errors:
				logger.FdebuglnFatal("File watcher error", err)
				logger.FatalLogln("File watcher error", err)
			}
		}
	}()

	// Don't wait for the first save event to send the code.
	// Send it as soon as user calls 'foundry go'
	initialUploadCh <- struct{}{}

	<-done
}

func ignored(s string, globs []glob.Glob) bool {
	logger.Fdebugln("string to match:", s)
	for _, g := range globs {
		logger.Fdebugln("\t- glob:", g)
		logger.Fdebugln("\t- match:", g.Match(s))
		if g.Match(s) {
			return true
		}
	}
	return false
}

func listenCallback(data []byte, err error) {
	logger.Fdebugln(string(data))

	if err != nil {
		logger.FdebuglnFatal("WebSocket error", err)
		logger.FatalLogln("WebSocket error", err)
	}

	t := connMsg.ResponseMsgType{}
	if err := json.Unmarshal(data, &t); err != nil {
		logger.FdebuglnFatal("Unmarshaling response error", err)
		logger.FatalLogln("Parsing server JSON response error", err)
	}

	_ = prompt.HideLoading()
	switch t.Type {
	case connMsg.LogResponseMsg:
		var s struct{ Content connMsg.LogContent }

		if err := json.Unmarshal(data, &s); err != nil {
			logger.FdebuglnFatal("Unmarshaling response error", err)
			logger.FatalLogln("Parsing server log message error", err)
		}

		if _, err := prompt.Writeln(s.Content.Msg); err != nil {
			logger.FdebuglnFatal("Error writing output", err)
			logger.FatalLogln("Error writing output", err)
		}

	case connMsg.WatchResponseMsg:
		var s struct{ Content connMsg.WatchContent }

		if err := json.Unmarshal(data, &s); err != nil {
			logger.FdebuglnFatal("Unmarshaling response error", err)
			logger.FatalLogln("Parsing server wathc message error", err)
		}

		var info string
		if s.Content.RunAll {
			info = "All filters disabled. Will display output from all functions."
		} else {
			info = fmt.Sprintf("Displaying output from: %s.", strings.Join(s.Content.Run, ", "))
		}

		prompt.SetInfoln(info, p.InfoLineSeverityWarning)
	}
}
