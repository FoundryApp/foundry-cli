package prompt

import (
	"fmt"
	"os"
	"strings"

	goprompt "github.com/mlejva/go-prompt"
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

func (c *Cmd) ToSuggest() goprompt.Suggest {
	return goprompt.Suggest{Text: c.Text, Description: c.Desc}
}

type Prompt struct {
	cmds 	[]*Cmd
}

var (
	promptText = ""
	promptTextCols = 0
	errorRow = 0
	promptRow = 0

	totalRows = 0
	freeRows = 0

	overlapping = false

	writer = goprompt.NewStandardOutputWriter()
)

func NewPrompt(cmds []*Cmd) *Prompt {
	return &Prompt{cmds}
}

func (p *Prompt) completer(d goprompt.Document) []goprompt.Suggest {
	s := []goprompt.Suggest{}
	for _, c := range p.cmds {
		s = append(s, c.ToSuggest())
	}
	return goprompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func (p *Prompt) executor(s string) {
	if s == "" { return }

	fields := strings.Fields(s)

	if cmd := p.getCommand(fields[0]); cmd != nil {
		args := fields[1:]

		if err := cmd.Do(args); err != nil {
			fmt.Println(err)
			// os.Exit(1)
		}
	} else {
		fmt.Printf("Unknown command '%s'. Write 'help' to list available commands.\n", fields[0])
	}
}

func (p *Prompt) getCommand(s string) *Cmd {
	for _, c := range p.cmds {
		if c.Text == s {
			return c
		}
	}
	return nil
}

func (p *Prompt) Print(t string) {

}

func (p *Prompt) Run() {
	parser := goprompt.NewStandardInputParser()
	size := parser.GetWinSize()

	totalRows = int(size.Row)
	promptRow = totalRows
	errorRow = promptRow - 1
	freeRows = promptRow - 3

	wReset()

	interup := goprompt.OptionAddKeyBind(goprompt.KeyBind{
		Key: 	goprompt.ControlC,
		Fn: 	func(buf *goprompt.Buffer) {
						os.Exit(0)
					},
	})
	newp := goprompt.New(p.executor, p.completer, interup)
	newp.Run()
}


func wReset() {
	writer.EraseScreen()
	writer.CursorGoTo(promptRow, 0)
	writer.Flush()
}

// TODO: Handle resizing terminal