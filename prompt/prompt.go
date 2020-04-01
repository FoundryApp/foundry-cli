package prompt

import (
	"fmt"
	"foundry/cli/logger"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"foundry/cli/prompt/cmd"

	goprompt "github.com/mlejva/go-prompt"
)

type CursorPos struct {
	Row int
	Col int
}

func CursorIdentity() CursorPos {
	return CursorPos{1, 1}
}

type PromptSafe struct {
	cmds []cmd.Cmd

	outBuf *Buffer
	// outBufMutex sync.Mutex

	renderMutex sync.Mutex

	promptPrefix string
	promptText   string
	promptRow    int // Will be recalculated once the terminal is ready

	errorText string
	errorRow  int // Will be recalculated once the terminal is ready

	totalRows int // Will be recalculated once the terminal is ready
	freeRows  int // Will be recalculated once the terminal is ready

	parser *goprompt.PosixParser
	writer goprompt.ConsoleWriter

	savedPos   CursorPos
	currentPos CursorPos // Current position of the cursor when printing output
}

//////////////////////

func (p *PromptSafe) completer(d goprompt.Document) []goprompt.Suggest {
	p.renderMutex.Lock()
	p.promptText = d.CurrentLine()
	p.renderMutex.Unlock()

	return []goprompt.Suggest{}
}

func (p *PromptSafe) executor(s string) {
	if s == "" {
		return
	}
	logger.Fdebugln(s)

	fields := strings.Fields(s)

	if cmd := p.getCommand(fields[0]); cmd != nil {
		logger.Fdebugln(cmd)
		args := fields[1:]
		cmd.RunRequest(args)
	} else {
		// Delete an old error message and show the new one

		p.renderMutex.Lock()

		// Delete an old error message
		p.writer.CursorGoTo(p.errorRow, 1)
		p.writer.EraseLine()

		// Print the new error
		p.writer.SetColor(goprompt.Red, goprompt.DefaultColor, true)
		p.errorText = fmt.Sprintf("Unknown command '%s'", fields[0])
		p.writer.WriteRawStr(p.errorText)
		p.writer.SetColor(goprompt.DefaultColor, goprompt.DefaultColor, false)

		// Move cursor back to the prompt
		p.writer.CursorGoTo(p.promptRow, len(p.promptPrefix)+len(p.promptText)+1)

		if err := p.writer.Flush(); err != nil {
			logger.Fdebugln(err)
			logger.ErrorLoglnFatal(err)
		}

		p.renderMutex.Unlock()
	}
}

func (p *PromptSafe) getCommand(s string) cmd.Cmd {
	for _, c := range p.cmds {
		if c.Name() == s {
			return c
		}
	}
	return nil
}

/////////////

func NewPromptSafe(cmds []cmd.Cmd) *PromptSafe {
	prefix := "> "
	return &PromptSafe{
		cmds: cmds,

		outBuf: NewBuffer(),

		promptPrefix: prefix,

		parser: goprompt.NewStandardInputParser(),
		writer: goprompt.NewStandardOutputWriter(),

		// Terminal is indexed from 1
		savedPos:   CursorIdentity(),
		currentPos: CursorPos{1, len(prefix) + 1},
	}
}

func (p *PromptSafe) Run() {
	// Read buffer and print anything that gets send to the channel
	bufCh := make(chan []byte, 128)
	stopReadCh := make(chan struct{})
	go p.outBuf.Read(bufCh, stopReadCh)
	go func() {
		for {
			select {
			case b := <-bufCh:
				p.print(b)
			default:
				time.Sleep(time.Millisecond * 10)
			}
		}
	}()

	interupOpt := goprompt.OptionAddKeyBind(goprompt.KeyBind{
		Key: goprompt.ControlC,
		Fn: func(buf *goprompt.Buffer) {
			os.Exit(0)
		},
	})
	prefixOpt := goprompt.OptionPrefix(p.promptPrefix)
	prefixColOpt := goprompt.OptionPrefixTextColor(goprompt.Green)
	prompt := goprompt.New(p.executor, p.completer, interupOpt, prefixOpt, prefixColOpt)
	go prompt.Run()

	// The initial rerender for the current terminal size
	if err := p.rerender(); err != nil {
		logger.Fdebugln(err)
		logger.ErrorLoglnFatal(err)
	}

	// Rerender a terminal for every size change
	go p.rerenderOnTermSizeChange()
}

func (p *PromptSafe) Writeln(s string) (n int, err error) {
	return p.outBuf.Write([]byte(s + "\n"))
}

func (p *PromptSafe) rerender() error {
	p.renderMutex.Lock()
	defer p.renderMutex.Unlock()

	p.writer.EraseScreen()

	p.currentPos = CursorIdentity()
	p.savedPos = CursorIdentity()

	size := p.parser.GetWinSize()
	p.totalRows = int(size.Row)
	p.promptRow = p.totalRows
	p.errorRow = p.totalRows - 1
	p.freeRows = p.totalRows

	// Move to the error row and restore the error text
	p.writer.CursorGoTo(p.errorRow, 1)
	p.writer.SetColor(goprompt.Red, goprompt.DefaultColor, true)
	p.writer.WriteRawStr(p.errorText)

	p.writer.CursorGoTo(p.promptRow, 1)

	return p.writer.Flush()
}

func (p *PromptSafe) rerenderOnTermSizeChange() {
	sigwinchCh := make(chan os.Signal, 1)
	defer close(sigwinchCh)
	signal.Notify(sigwinchCh, syscall.SIGWINCH)
	for {
		if _, ok := <-sigwinchCh; !ok {
			return
		}
		if err := p.rerender(); err != nil {
			logger.FdebuglnFatal(err)
			logger.ErrorLoglnFatal(err)
		}
	}
}

func (p *PromptSafe) print(b []byte) {
	p.renderMutex.Lock()
	defer p.renderMutex.Unlock()

	// The invariant is that the the p.savedPos always holds
	// a position where we stopped printing the text = where
	// we should start printing text again.
	p.writer.CursorGoTo(p.savedPos.Row, p.savedPos.Col)

	s := string(b)
	logger.Fdebugln(s)
	for _, r := range s {
		p.writer.WriteRawStr(string(r))
		p.currentPos.Col++

		if r == '\n' {
			// On a new line, the cursor moves to the start of a line
			p.currentPos.Col = 1

			p.currentPos.Row++
			p.freeRows--
		}

		if p.freeRows == 2 {
			p.savedPos = p.currentPos
			// Go to a prompt row and create a new line so that we
			// once again have 3 free rows.
			// The reason we have to go to the prompt row is becauase
			// if we had printed a new line anywhere before the prompt
			// row, the cursor would simply move down without actually
			// creating a new line in the terminal.

			// Erase the error row and prompt row  so that a text doesn't stay there
			// when the everything is moved up by 1 row
			p.writer.CursorGoTo(p.errorRow, 1)
			p.writer.EraseLine()
			p.writer.CursorGoTo(p.promptRow, 1)
			p.writer.EraseLine()

			// Create a new line
			p.writer.WriteRawStr("\n")

			// Move cursor back to a position where we stopped outputting
			// text. This will be next available new line after the last
			// line of printed text
			p.writer.CursorGoTo(p.savedPos.Row, p.savedPos.Col)
			// The reason it's not sufficient to just go to p.savedPos
			// is because we printed a newline. All text moved 1 line up.
			p.writer.CursorUp(1)

			p.currentPos.Row--
			p.currentPos.Col = 1
			p.freeRows = 3
		}
	}
	p.savedPos = p.currentPos

	// Move to the error row and restore the error text
	p.writer.CursorGoTo(p.errorRow, 1)
	p.writer.SetColor(goprompt.Red, goprompt.DefaultColor, true)
	p.writer.WriteRawStr(p.errorText)

	// Move to the prompt row and restore the text
	p.writer.CursorGoTo(p.promptRow, 1)
	p.writer.SetColor(goprompt.Green, goprompt.DefaultColor, false)
	p.writer.WriteRawStr(p.promptPrefix)
	p.writer.SetColor(goprompt.DefaultColor, goprompt.DefaultColor, false)
	p.writer.WriteRawStr(p.promptText)

	if err := p.writer.Flush(); err != nil {
		logger.Fdebugln(err)
		logger.ErrorLoglnFatal(err)
	}
}
