package cmd

import (
	"fmt"
	c "foundry/cli/connection"
	connMsg "foundry/cli/connection/msg"

	goprompt "github.com/mlejva/go-prompt"
)

type WatchCmd struct {
	Text  string
	Desc  string
	RunCh RunChannelType
}

func NewWatchCmd() *WatchCmd {
	return &WatchCmd{
		Text:  "watch",
		Desc:  "Watch only specific function(s)",
		RunCh: make(chan Args),
	}
}

func NewWatchAllCmd() *WatchCmd {
	return &WatchCmd{
		Text:  "watch:all",
		Desc:  "Disable all active watch filters and watch all functions",
		RunCh: make(chan Args),
	}
}

// Implement Cmd interface
func (c *WatchCmd) Run(conn *c.Connection, args Args) (promptOutput string, promptInfo string, err error) {
	watchAll := false
	fns := args
	if c.Text == "watch:all" {
		watchAll = true
		fns = []string{}
	} else {
		if len(args) == 0 {
			return "", "No argument specified. Example usage: 'watch myFunction'", nil
		}
	}

	msg := connMsg.NewWatchfnMsg(watchAll, fns)
	err = conn.Send(msg)
	return "", "", err
}

func (c *WatchCmd) RunRequest(args Args) {
	c.RunCh <- args
}

func (c *WatchCmd) ToSuggest() goprompt.Suggest {
	return goprompt.Suggest{c.Text, c.Desc}
}

func (c *WatchCmd) Name() string {
	return c.Text
}

func (c *WatchCmd) String() string {
	return fmt.Sprintf("%s - %s", c.Text, c.Desc)
}
