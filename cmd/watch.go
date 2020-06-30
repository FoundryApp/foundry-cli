package cmd

import (
	"foundry/cli/desktopapp"
	"foundry/cli/logger"
	cliprompt "foundry/cli/prompt"
	promptcmd "foundry/cli/prompt/cmd"
	"foundry/cli/rwatch"
	"foundry/cli/session"
	"foundry/cli/user"
	"foundry/cli/zip"
	"os"

	"github.com/spf13/cobra"
)

var (
	watchCmd = &cobra.Command{
		Use:     "watch",
		Short:   "Foundry starts watching code files for changes. With every change, Foundry automaticaly hot-reloads code in the deployed service.",
		Example: "foundry watch",
		Run:     runWatch,
	}

	exitCmd *promptcmd.ExitCmd
	prompt  *cliprompt.Prompt
	sess    *session.Session
)

func init() {
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) {
	done := make(chan struct{})

	CheckForPackageJSON()
	LoadIgnoreFile()

	checkAppStatus()
	// Request the current user and session from the desktop app
	user, err := user.GetCurrent()
	if err != nil {
		logger.DebuglnError("Error while getting current user from the app", err)
		logger.FatalLogln("Couldn't get info about the current user", err)
	}

	s, err := session.GetCurrent()
	if err != nil {
		logger.DebuglnError("Error while getting current session from the app", err)
		logger.FatalLogln("Couldn't get info about the current session", err)
	}
	sess = s

	if err := sess.Connect(user.Creds.Token); err != nil {
		logger.DebuglnError("Error connection to the session", err)
		logger.FatalLogln("Couldn't connect to the environment", err)
	}

	// File watcher
	fwatcher, err := setUpFileWatcher()
	if err != nil {
		logger.FatalLogln(err)
	}
	defer fwatcher.Close()
	go fwatcher.Watch()

	// Session
	data := make(chan []byte)
	listenErr := make(chan error)
	go sess.Listen(data, listenErr)
	go handleListenChannels(data, listenErr)

	// Prompt
	prompt = setUpPrompt()
	go prompt.Run()

	// The goroutine handling all file events + prompt command requests.
	// Command requests are all handled from a single goroutine because
	// Gorilla's websocket connection supports only one concurrent reader
	// and one concurrent writer.
	// More info - https://godoc.org/github.com/gorilla/websocket#hdr-Concurrency
	go func() {
		for {
			select {
			case e := <-prompt.Events:
				handlePromptEvent(&e)
			case /*cmdArgs*/ _ = <-exitCmd.RunCh:
				// TODO: RUN COMMAND
				// _, _, _ = exitCmd.Run(connectionClient, args)
			case e := <-fwatcher.Events:
				path := "." + string(os.PathSeparator) + e.Name
				if !ignored(path, foundryConf.Ignore) {
					logger.Fdebugln("Watcher event", e.Name)
					_ = prompt.ShowLoading()
					zipAndSend()
				}
			case err := <-fwatcher.Errors:
				logger.FdebuglnFatal("File watcher error", err)
				logger.FatalLogln("Error - ", err)
			}
		}
	}()

	// Never finish watch command executionso the prompt doesn't exit
	<-done
}

func handlePromptEvent(event *cliprompt.PromptEvent) {
	switch event.Type {
	case cliprompt.PromptEventTypeRerender:
		if err := prompt.ShowLoading(); err != nil {
			// TODO: Handle error
		}
		zipAndSend()
	}
}

func zipAndSend() error {
	buf, err := zip.ArchiveDir(foundryConf.CurrentDir, foundryConf.Ignore)
	if err != nil {
		logger.FdebuglnError("ArchiveDir error:", err)
		return err
	}
	return sess.SendData(buf)
}

func setUpFileWatcher() (*rwatch.Watcher, error) {
	fwatcher, err := rwatch.New(foundryConf.Ignore)
	if err != nil {
		logger.FdebuglnFatal("Watcher initialization error: ", err)
		return nil, err
	}

	if err = fwatcher.AddRecursive(foundryConf.CurrentDir); err != nil {
		logger.FdebuglnFatal("Watcher AddRecursive error: ", err)
		return nil, err
	}
	return fwatcher, nil
}

func setUpPrompt() *cliprompt.Prompt {
	exitCmd = promptcmd.NewExitCmd()
	cmds := []promptcmd.Cmd{exitCmd}
	return cliprompt.NewPrompt(cmds)
}

func handleListenChannels(data chan []byte, err chan error) {
	// TODO
	for {
		select {
		case d := <-data:
			logger.Debugln(string(d))
		case e := <-err:
			logger.DebuglnError("WebSocket error: ", e)
			logger.FatalLogln("WebSocket error: ", e)
		}
	}
}

func checkAppStatus() {
	status, err := desktopapp.GetStatus()
	if err != nil {
		logger.DebuglnError("Couldn't get status from the Foundry Desktop app: ", err)
		logger.FatalLogln("Couldn't connect to Foundry Desktop App. Is the app running? Download it at https://foundryapp.co/download")
	}
	switch status {
	case desktopapp.ReadyAppStatus:
		// Ok
	case desktopapp.NotReadyAppStatus:
		logger.FatalLogln("Foundry Desktop App isn't fully loaded yet, please wait a few seconds and try again")
	default:
		logger.DebuglnError("Unknown desktop app status: ", status)
		logger.FatalLogln("Couldn't correctly communiciate with the Foundry Desktop App")
	}
}
