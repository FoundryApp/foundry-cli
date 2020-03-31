package prompt

import (
	"bytes"
	"foundry/cli/logger"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	goprompt "github.com/mlejva/go-prompt"
)

type SavedPos struct {
	row int
	col int
}

type PromptSafe struct {
	cmds []*Cmd

	outBuf      bytes.Buffer
	outBufMutex sync.Mutex

	renderMutex sync.Mutex

	promptPrefix string
	promptText   string
	promptRow    int

	errorText string
	errorRow  int

	totalRows int
	freeRows  int

	parser *goprompt.PosixParser
	writer goprompt.ConsoleWriter

	waitDuration time.Duration

	savedPos SavedPos
	col      int
	row      int
}

//////////////////////

func (p *PromptSafe) executor(s string) {
}

func (p *PromptSafe) completer(d goprompt.Document) []goprompt.Suggest {
	p.promptText = d.CurrentLine()
	return []goprompt.Suggest{}
}

/////////////

func NewPromptSafe() *PromptSafe {
	return &PromptSafe{
		outBuf: bytes.Buffer{},

		promptPrefix: "> ",
		promptText:   "",
		promptRow:    0,

		errorText: "",
		errorRow:  0, // Will be recalculated once the terminal is ready

		totalRows: 0, // Will be recalculated once the terminal is ready
		freeRows:  0, // Will be recalculated once the terminal is ready

		parser: goprompt.NewStandardInputParser(),
		writer: goprompt.NewStandardOutputWriter(),

		// Terminal is indexed from 1
		savedPos: SavedPos{1, 1},
		row:      1,
		col:      len("> ") + 1,
		// waitDuration: time.Microsecond * 400,
	}
}

func (p *PromptSafe) WriteOutputln(s string) (n int, err error) {
	p.outBufMutex.Lock()
	defer p.outBufMutex.Unlock()
	return p.outBuf.Write([]byte(s + "\n"))
}

func (p *PromptSafe) Run() {
	bufCh := make(chan []byte, 128)
	stopReadCh := make(chan struct{})

	if err := p.rerender(); err != nil {
		logger.Fdebugln(err)
		logger.LoglnFatal(err)
	}

	// Watch for terminal size changes
	sigwinch := make(chan os.Signal, 1)
	defer close(sigwinch)
	signal.Notify(sigwinch, syscall.SIGWINCH)
	go func() {
		for {
			if _, ok := <-sigwinch; !ok {
				return
			}
			if err := p.rerender(); err != nil {
				logger.Fdebugln(err)
				logger.LoglnFatal(err)
			}
		}
	}()

	// Read buffer and print anything that gets send to the channel
	go p.readOutBuffer(bufCh, stopReadCh)
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
	prompt := goprompt.New(p.executor, p.completer, interupOpt)
	prompt.Run()
}

func (p *PromptSafe) readOutBuffer(bufCh chan<- []byte, stopCh <-chan struct{}) {
	for {
		select {
		case <-stopCh:
			return
		default:
			p.outBufMutex.Lock()

			buf := make([]byte, 1024)
			n, err := p.outBuf.Read(buf)

			if err == nil {
				bufCh <- buf[:n]
			} else if err != io.EOF {
				logger.Fdebugln(err)
				logger.LoglnFatal(err)
			}

			p.outBufMutex.Unlock()
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func (p *PromptSafe) rerender() error {
	p.renderMutex.Lock()
	defer p.renderMutex.Unlock()

	p.row = 1
	p.col = 1
	p.savedPos = SavedPos{1, 1}

	size := p.parser.GetWinSize()
	p.totalRows = int(size.Row)
	p.promptRow = p.totalRows
	p.errorRow = p.totalRows - 1
	p.freeRows = p.totalRows

	writer.EraseScreen()

	// TODO: Restore error

	// Restore prompt
	writer.CursorGoTo(p.promptRow, 1)
	writer.WriteRawStr(p.promptPrefix + p.promptText)

	return writer.Flush()
}

func (p *PromptSafe) print(b []byte) {
	p.renderMutex.Lock()
	defer p.renderMutex.Unlock()

	writer.CursorGoTo(p.savedPos.row, p.savedPos.col)
	if err := writer.Flush(); err != nil {
		logger.Fdebugln(err)
		logger.LoglnFatal(err)
	}

	s := string(b)
	logger.Fdebugln(s)

	for _, r := range s {
		writer.WriteRawStr(string(r))
		if err := writer.Flush(); err != nil {
			logger.Fdebugln(err)
			logger.LoglnFatal(err)
		}

		p.col++
		if r == '\n' {
			p.col = 1
			p.row++

			p.freeRows--
		}

		if p.freeRows == 2 {
			p.savedPos = SavedPos{p.row, p.col}
			if err := writer.Flush(); err != nil {
				logger.Fdebugln(err)
				logger.LoglnFatal(err)
			}

			// TODO: Erase error

			writer.CursorGoTo(p.promptRow, 1)
			writer.EraseLine()
			if err := writer.Flush(); err != nil {
				logger.Fdebugln(err)
				logger.LoglnFatal(err)
			}

			writer.WriteRawStr("\n")
			if err := writer.Flush(); err != nil {
				logger.Fdebugln(err)
				logger.LoglnFatal(err)
			}

			p.row--
			p.col = 1
			p.freeRows = 3

			// Restore prompt
			writer.CursorGoTo(p.promptRow, 1)
			writer.WriteRawStr(p.promptPrefix + p.promptText)
			if err := writer.Flush(); err != nil {
				logger.Fdebugln(err)
				logger.LoglnFatal(err)
			}

			// Go back to the next available output line
			writer.CursorGoTo(p.savedPos.row, p.savedPos.col)
			if err := writer.Flush(); err != nil {
				logger.Fdebugln(err)
				logger.LoglnFatal(err)
			}
			writer.CursorUp(1)
			if err := writer.Flush(); err != nil {
				logger.Fdebugln(err)
				logger.LoglnFatal(err)
			}
		}
	}
	p.savedPos = SavedPos{p.row, p.col}

	writer.CursorGoTo(p.promptRow, len(p.promptPrefix)+1+len(p.promptText)+1)
	if err := writer.Flush(); err != nil {
		logger.Fdebugln(err)
		logger.LoglnFatal(err)
	}

	// TODO: Restore error

}
