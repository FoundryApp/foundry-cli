package prompt

import (
	// "bytes"
	"fmt"
	// "io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"foundry/cli/logger"

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

	// buf *bytes.Buffer
}

var (
	promptPrefix = "> "

	promptText = ""
	promptRow = 0

	errorText = ""
	errorRow = 0

	totalRows = 0
	freeRows = 0

	overlapping = false
	overlappingRows = 0

	parser = goprompt.NewStandardInputParser()
	writer = goprompt.NewStandardOutputWriter()

	wsaved = false

	waitDuration = time.Millisecond * 300
)

func NewPrompt(cmds []*Cmd) *Prompt {
	return &Prompt{cmds,/* &bytes.Buffer{}*/}
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
			logger.FdebuglnFatal(err)
			logger.LogFatal(err)
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

// func (p *Prompt) WriteToBuffer(s string) error {
// 	_, err := p.buf.Write([]byte(s))
// 	return err
// }

// func (p *Prompt) watchBuffer() {
// 	for {
// 		// logger.Fdebugln("Watch Buffer")

// 		b := make([]byte, 1024)
// 		if n, err := p.buf.Read(b); err == nil && n > 0 {
// 			p.Print(string(b))
// 			// logger.Fdebugln("BUFFER:", string(b))
// 		} else if err != nil && err != io.EOF {
// 			logger.FdebuglnFatal(err)
// 			logger.LogFatal(err)
// 		}

// 		// time.Sleep(time.Millisecond * 10)
// 	}
// }

func (p *Prompt) Print(s string) {
	logger.Fdebugln("#### START PRINT")

	logger.Fdebugln("[print] totalRows:", totalRows)
	logger.Fdebugln("[print] promptRow:", promptRow)
	logger.Fdebugln("[print] errorRow:", errorRow)

	logger.Fdebugln("[print] raw:", s)
	trimmed := strings.TrimSpace(s)
	logger.Fdebugln("[print] trimmed:", trimmed)
	lines := strings.Split(trimmed, "\n")
	logger.Fdebugln("[print] totalLines:", len(lines))

	for _, l := range lines {
		logger.Fdebugln("[prompt] freeRows start:", freeRows)
		logger.Fdebugln("[prompt] line:", l)

		freeRows -= 1

		writer.UnSaveCursor()
		writer.Flush()

		writer.WriteRawStr(l+"\n")
		writer.Flush()

		writer.SaveCursor()
		writer.Flush()

		if freeRows <= 3 {
			newRows := 3 - freeRows
			logger.Fdebugln("[prompt] newRows:", newRows)

			writer.CursorGoTo(errorRow, 0)
			writer.Flush()
			writer.EraseLine()
			writer.Flush()

			writer.CursorGoTo(promptRow, 0)
			writer.Flush()
			writer.EraseLine()
			writer.Flush()
			writer.WriteRawStr(strings.Repeat("\n", newRows))
			writer.Flush()

			freeRows += newRows

			writer.UnSaveCursor()
			writer.Flush()

			writer.CursorUp(newRows)
			writer.Flush()
			writer.SaveCursor()
			writer.Flush()
		}

		logger.Fdebugln("[prompt] freeRows end:", freeRows)
	}

	writer.CursorGoTo(errorRow, 0)
	writer.Flush()
	writer.WriteRawStr(errorText)
	writer.Flush()

	writer.CursorGoTo(promptRow, 0)
	writer.Flush()
	writer.WriteRawStr(promptPrefix + promptText)
	writer.Flush()

	logger.Fdebugln("#### END PRINT")
}

func (p *Prompt) SetPromptPrefix(s string) {
	promptPrefix = s
}

func (p *Prompt) Run() {
	size := parser.GetWinSize()

	// totalRows = int(size.Row)
	// promptRow = totalRows
	// errorRow = promptRow - 1
	// freeRows = totalRows



	// Watch for terminal size changes
	sigwinch := make(chan os.Signal, 1)
	defer close(sigwinch)
	signal.Notify(sigwinch, syscall.SIGWINCH)
	go func() {
		for {
			if _, ok := <-sigwinch; !ok { return }
			size = parser.GetWinSize()
			logger.Fdebugln("Terminal size change:", size)
			p.rerender(size)
		}
	}()

	p.rerender(size)

	// writer.CursorGoTo(0, 0)
	// writer.Flush()
	// writer.SaveCursor()
	// writer.Flush()

	interupOpt := goprompt.OptionAddKeyBind(goprompt.KeyBind{
		Key: 	goprompt.ControlC,
		Fn: 	func(buf *goprompt.Buffer) {
						os.Exit(0)
					},
	})
	prefixOpt := goprompt.OptionPrefix(promptPrefix)
	livePrefixOpt := goprompt.OptionLivePrefix(func() (prefix string, useLivePrefix bool) {
		return promptPrefix, true
	})

	newp := goprompt.New(p.executor, p.completer, interupOpt, prefixOpt, livePrefixOpt)

	// go p.watchBuffer()

	newp.Run()
}

// func (p *Prompt) calcOverlapping(t string) {
// 	l := strings.Split(t, "\n")

// 	if len(l) >= freeRows {
// 		freeRows = 0
// 		overlapping = true
// 		overlappingRows = int(math.Abs(float64(freeRows - len(l))))
// 	} else {
// 		freeRows -= len(l)
// 		overlappingRows = 0
// 	}
// }

func (p *Prompt) rerender(size *goprompt.WinSize) {
	totalRows = int(size.Row)
	promptRow = totalRows
	errorRow = promptRow - 1
	freeRows = totalRows

	// So the initial UnSave is at 0,0
	writer.CursorGoTo(0, 0)
	writer.Flush()
	writer.SaveCursor()
	writer.Flush()

	// Clears the screen and moves cursor to promptRow
	p.wReset()

	logger.Fdebugln("totalRows:", totalRows)
	logger.Fdebugln("promptRow:", promptRow)
	logger.Fdebugln("errorRow:", errorRow)
	logger.Fdebugln("freeRows:", freeRows)
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