package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	connMsg "foundry/cli/connection/msg"
	"foundry/cli/desktopapp"
	"foundry/cli/logger"
	"foundry/cli/rwatch"
	"foundry/cli/session"
	"foundry/cli/user"
	"foundry/cli/zip"

	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
)

var (
	watchCmd = &cobra.Command{
		Use:     "watch",
		Short:   "Foundry starts watching code files for changes. With every change, Foundry automaticaly hot-reloads code in the deployed service.",
		Example: "foundry watch",
		Run:     runWatch,
	}

	sess *session.Session
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
	sessData := make(chan []byte)
	sessErr := make(chan error)
	go sess.Listen(sessData, sessErr)

	triggerInitialSend := make(chan struct{}, 1)

	// The goroutine handling all file events + prompt command requests.
	// Command requests are all handled from a single goroutine because
	// Gorilla's websocket connection supports only one concurrent reader
	// and one concurrent writer.
	// More info - https://godoc.org/github.com/gorilla/websocket#hdr-Concurrency
	go func() {
		for {
			select {
			case d := <-sessData:
				logger.Fdebugln(string(d))
				parseMessageData(d)
			case e := <-sessErr:
				logger.DebuglnError("Session error: ", e)
				logger.FatalLogln("Session error: ", e)
			case e := <-fwatcher.Events:
				path := "." + string(os.PathSeparator) + e.Name
				if !ignored(path, foundryConf.Ignore) {
					zipAndSend()
				}
			case err := <-fwatcher.Errors:
				logger.FdebuglnFatal("File watcher error: ", err)
				logger.FatalLogln("Error: ", err)
			case _ = <-triggerInitialSend:
				zipAndSend()
			}
		}
	}()

	// Don't wait for the first save event to send the code.
	// Send it as soon as user calls 'foundry watch'
	triggerInitialSend <- struct{}{}

	// Never finish watch command execution so the log stream doesn't exit
	<-done
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
func parseMessageData(d []byte) {
	msg := connMsg.ResponseMsgType{}
	if err := json.Unmarshal(d, &msg); err != nil {
		logger.DebuglnError("Couldn't parse environment response: ", err)
		logger.FatalLogln("Couldn't parse environment response: ", err)
	}

	switch msg.Type {
	case connMsg.LogResponseMsg:
		var s struct{ Content connMsg.LogContent }
		if err := json.Unmarshal(d, &s); err != nil {
			logger.DebuglnError("Couldn't parse environment log message: ", err)
			logger.FatalLogln("Couldn't parse environment log message: ", err)
		}
		fmt.Print(s.Content.Msg)
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
