package cmd

import (
	"fmt"
	c "foundry/cli/connection"
	connMsg "foundry/cli/connection/msg"
	"foundry/cli/logger"

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

// Implement Cmd interface

func (c *WatchCmd) Run(conn *c.Connection, args Args) error {
	if len(args) == 0 {
		logger.Logln("Write 'watch all' to watch all functions or 'watch <function-name1> <function-name2> ...' to watch specific functions")
		return nil
	}

	watchAll := false
	fns := args
	if args[0] == "all" {
		watchAll = true
		fns = []string{}
	}

	msg := connMsg.NewWatchfnMsg(watchAll, fns)
	return conn.Send(msg)
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
