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
	// TODO: vars should be here? At least writer

	// buff bytes.Buffer
}

var (
	promptText = ""
	promptRow = 0

	errorText = ""
	errorRow = 0

	totalRows = 0
	freeRows = 0

	overlapping = false

	writer = goprompt.NewStandardOutputWriter()

	wsaved = false

	f *os.File
)

func NewPrompt(cmds []*Cmd) *Prompt {
	return &Prompt{cmds}
}

func (p *Prompt) completer(d goprompt.Document) []goprompt.Suggest {
	promptText = d.CurrentLine()

	s := []goprompt.Suggest{}
	for _, c := range p.cmds {
		s = append(s, c.ToSuggest())
	}

	return []goprompt.Suggest{}
	//return goprompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
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
		p.wGoToAndEraseError()

		errorText = fmt.Sprintf("Unknown command '%s'. Write 'help' to list available commands.\n", fields[0])
		writer.WriteStr(errorText)
		writer.Flush()

		p.wGoToPrompt()
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
	// TODO: Handle resizing terminal
	// TODO: Handle text that is too long and is rendered as a multiline text

	t = strings.TrimSpace(t)
	fmt.Fprintln(f, t)


	// s := fmt.Sprintf("split: %v", split)
	// fmt.Fprintln(f, s)

	split := strings.Split(t, "\n")
	l := len(split)
	s := fmt.Sprintf("LEN: %v", l)
	fmt.Fprintln(f, s)

	s = fmt.Sprintf("wsaved: %v", wsaved)
	fmt.Fprintln(f, s)

	s = fmt.Sprintf("freeRows: %v", freeRows)
	fmt.Fprintln(f, s)

	s = fmt.Sprintf("overlapping: %v", overlapping)
	fmt.Fprintln(f, s)

	if wsaved {
		writer.UnSaveCursor()
		writer.Flush()
		if overlapping {
			// writer.CursorUp(2)
		}
		wsaved = false
	} else {
		writer.CursorGoTo(0, 0)
	}

	writer.SaveCursor()
	writer.Flush()

	p.wGoToAndErasePrompt()
	p.wGoToAndEraseError()

	// Restore the cursor
	writer.UnSaveCursor()
	writer.Flush()

	// Output the text
	p.calcOverlapping(t)
	writer.WriteRawStr(t+"\n")
	writer.Flush()
	writer.SaveCursor()
	wsaved = true
	writer.Flush()

	// Create space for the prompt line + error line
	writer.WriteRawStr("\n")
	writer.Flush()

	// Restore the error
	p.wGoToError()
	writer.WriteStr(errorText)
	writer.Flush()

	// Restore the prompt lines
	p.wGoToPrompt()
	writer.WriteStr("> " + promptText)
	writer.Flush()
}


func (p *Prompt) Run() {
	file, _ := os.Create("/Users/vasekmlejnsky/Developer/foundry/cli/debug.txt")
	f = file
	defer f.Close()

	parser := goprompt.NewStandardInputParser()
	size := parser.GetWinSize()

	totalRows = int(size.Row)
	promptRow = totalRows
	errorRow = promptRow - 1
	freeRows = promptRow - 3

	s := fmt.Sprintf("totalRows: %v", totalRows)
	fmt.Fprintln(f, s)
	s = fmt.Sprintf("promptRow: %v", promptRow)
	fmt.Fprintln(f, s)
	s = fmt.Sprintf("errorRow: %v", errorRow)
	fmt.Fprintln(f, s)
	s = fmt.Sprintf("freeRows: %v", freeRows)
	fmt.Fprintln(f, s)

	fmt.Fprintln(f, "=================")

	p.wReset()

	interup := goprompt.OptionAddKeyBind(goprompt.KeyBind{
		Key: 	goprompt.ControlC,
		Fn: 	func(buf *goprompt.Buffer) {
						os.Exit(0)
					},
	})
	newp := goprompt.New(p.executor, p.completer, interup)
	newp.Run()
}

func (p *Prompt) calcOverlapping(t string) {
	l := strings.Split(t, "\n")

	if len(l) >= freeRows {
		freeRows = 0
		overlapping = true
	} else {
		freeRows -= len(l)
	}
}

func (p *Prompt) wReset() {
	writer.EraseScreen()
	writer.CursorGoTo(promptRow, 0)
	writer.Flush()
}

func (p *Prompt) wGoToPrompt() {
	writer.CursorGoTo(promptRow, 0)
	writer.Flush()
}

func (p *Prompt) wGoToError() {
	writer.CursorGoTo(errorRow, 0)
	writer.Flush()
}

func (p *Prompt) wGoToAndErasePrompt() {
	p.wGoToPrompt()
	writer.EraseLine()
	writer.Flush()
}

func (p *Prompt) wGoToAndEraseError() {
	p.wGoToError()
	writer.EraseLine()
	writer.Flush()
}