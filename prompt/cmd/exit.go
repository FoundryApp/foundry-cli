package cmd

import (
	"fmt"
	c "foundry/cli/connection"
	"os"

	goprompt "github.com/mlejva/go-prompt"
)

type ExitCmd struct {
	Text  string
	Desc  string
	RunCh RunChannelType
}

func NewExitCmd() *ExitCmd {
	return &ExitCmd{
		Text:  "exit",
		Desc:  "Stop Foundry CLI",
		RunCh: make(chan Args),
	}
}

// Implement Cmd interface

func (c *ExitCmd) Run(conn *c.Connection, args Args) (promptOutput string, promptInfo string, err error) {
	os.Exit(0)
	return "", "", err
}

func (c *ExitCmd) RunRequest(args Args) {
	c.RunCh <- args
}

func (c *ExitCmd) ToSuggest() goprompt.Suggest {
	return goprompt.Suggest{c.Text, c.Desc}
}

func (c *ExitCmd) Name() string {
	return c.Text
}

func (c *ExitCmd) String() string {
	return fmt.Sprintf("%s - %s", c.Text, c.Desc)
}
