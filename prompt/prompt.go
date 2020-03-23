package prompt

import (
	"fmt"
	"github.com/c-bata/go-prompt"
)

type CmdRunFunc func(args []string) error

type Cmd struct {
	Text 	string
	Desc 	string
	Do		CmdRunFunc
}

func (c *Cmd) String() string {
	return fmt.Sprintf("%s - %s\n", c.Text, c.Desc)
}

func (c *Cmd) ToSuggest() prompt.Suggest {
	return prompt.Suggest{Text: c.Text, Description: c.Desc}
}
