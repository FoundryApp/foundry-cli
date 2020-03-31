package cmd

// "foundry go" or "foundry connect" or "foundry " or "foundry start" or "foundry link"?

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"foundry/cli/auth"
	conn "foundry/cli/connection"
	connMsg "foundry/cli/connection/msg"
	"foundry/cli/files"
	"foundry/cli/logger"
	p "foundry/cli/prompt"
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
	start = time.Now()

	// prompt      *p.Prompt
	prompt      *p.PromptSafe
	uploadStart = time.Now()

	df *os.File
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
			logger.Fdebugln(err)
			logger.ErrorLoglnFatal(err)
		}

		if authClient.Error != nil {
			logger.Fdebugln(authClient.Error)
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
		logger.Fdebugln("Connection error", err)
		logger.ErrorLoglnFatal(err)
	}
	defer c.Close()

	// Start an interactive prompt
	// cmds := []*p.Cmd{
	// 	// pc.Watch(c),
	// 	pc.Exit(),
	// }
	// prompt = p.NewPrompt(cmds)
	// go prompt.Run()

	// time.Sleep(time.Second * 20)

	prompt = p.NewPromptSafe()
	go prompt.Run()

	// Listen for messages from the WS connection
	go c.Listen(listenCallback)

	// Start periodically pinging server so the env isn't killed
	pingMsg := connMsg.NewPingMsg(conn.PingURL(), authClient.IDToken)
	ticker := time.NewTicker(time.Second * 10)
	go c.Ping(pingMsg, ticker, done)

	// Start the file watcher
	w, err := rwatch.New()
	if err != nil {
		logger.FdebuglnFatal("Watcher error", err)
		logger.LoglnFatal(err)
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
				logger.Fdebugln("watcher error", err)
				logger.ErrorLoglnFatal(err)
			}
		}
	}()

	err = w.AddRecursive(conf.RootDir)
	if err != nil {
		logger.FdebuglnFatal("watcher AddRecursive", err)
		logger.LoglnFatal(err)
	}

	// Don't wait for first save to send the code - send it as soon
	// as user calls 'foundry go'
	// logger.Fdebugln("<timer> reseting starting upload time")
	// uploadStart = time.Now()
	// files.Upload(c, conf.RootDir)

	<-done
}

func listenCallback(data []byte, err error) {
	// elapsed := time.Since(uploadStart)
	// logger.Fdebugln("<timer> time until response (+ time.Sleep) -", elapsed)

	// time.Sleep(time.Millisecond * 20)
	logger.Fdebugln(string(data))

	if err != nil {
		// elapsed := time.Since(start)
		// logger.Fdebugln("<timer> Elapsed time -", elapsed)
		logger.FdebuglnFatal("WS error", err)
		logger.LoglnFatal(err)
	}

	t := connMsg.ResponseMsgType{}
	if err := json.Unmarshal(data, &t); err != nil {
		// elapsed := time.Since(start)
		// logger.Fdebugln("<timer> Elapsed time -", elapsed)
		logger.FdebuglnFatal("Unmarshaling response error", err)
		logger.LoglnFatal(err)
	}

	switch t.Type {
	case connMsg.LogResponseMsg:
		var s struct{ Content connMsg.LogContent }

		if err := json.Unmarshal(data, &s); err != nil {
			// elapsed := time.Since(start)
			// logger.Fdebugln("<timer> Elapsed time -", elapsed)
			logger.FdebuglnFatal("Unmarshaling response error", err)
		}

		// logger.Logln(string(s.Content.Msg))

		// TODO: listenCallback is a callback - it doesn't wait for prompt to print everything
		// prompt must have a buffer and a lock that makes sure that it's printing sequentially

		s1 := fmt.Sprintf("[0] %s", s.Content.Msg)
		// s2 := fmt.Sprintf("[1] %s\n", s.Content.Msg)
		// s3 := fmt.Sprintf("[2] %s\n", s.Content.Msg)
		// s4 := fmt.Sprintf("[3] %s\n", s.Content.Msg)
		// s5 := fmt.Sprintf("[4] %s\n", s.Content.Msg)

		if _, err := prompt.WriteOutputln(s1); err != nil {
			logger.FdebuglnFatal(err)
		}

		// if _, err := prompt.WriteOutputln(s2); err != nil {
		// 	logger.FdebuglnFatal(err)
		// }

		// prompt.Write([]byte(s1))
		// prompt.Write([]byte(s2))
		// prompt.Write([]byte(s3))
		// prompt.Write([]byte(s4))
		// prompt.Write([]byte(s5))

		// prompt.Print(0, string(s.Content.Msg))
		// prompt.Print(1, string(s.Content.Msg))

	case connMsg.WatchResponseMsg:
		var s struct{ Content connMsg.WatchContent }

		if err := json.Unmarshal(data, &s); err != nil {
			// elapsed := time.Since(start)
			// logger.Fdebugln("<timer> Elapsed time -", elapsed)
			logger.FdebuglnFatal("Unmarshaling response error", err)
			logger.LoglnFatal(err)
		}

		// p := fmt.Sprintf("[%s] > ", strings.Join(s.Content.Run, ", "))
		// prompt.SetPromptPrefix(p)

		// p = fmt.Sprintf("Watching only functions: %s", strings.Join(s.Content.Run, ", "))
		// prompt.PrintInfo(p)
	}
}
