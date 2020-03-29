package cmd

import (
	c "foundry/cli/connection"
	connMsg "foundry/cli/connection/msg"
	"foundry/cli/prompt"
	"foundry/cli/logger"
)

var (
	conn *c.Connection
)

func Watch(cn *c.Connection) *prompt.Cmd {
	conn = cn
	return &prompt.Cmd{"watch", "Watch only specific function(s)", runWatch}
}

func runWatch(args []string) error {
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