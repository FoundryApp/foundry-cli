package cmd

import (
	"fmt"
	c "foundry/cli/connection"

	goprompt "github.com/mlejva/go-prompt"
)

type Args []string
type RunChannelType chan Args

type Cmd interface {
	Run(conn *c.Connection, args Args) (promptOutput string, promptInfo string, err error)
	RunRequest(args Args)
	ToSuggest() goprompt.Suggest
	Name() string
	fmt.Stringer
}
