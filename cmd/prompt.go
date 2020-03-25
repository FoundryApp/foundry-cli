package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"log"
	"sync"

	fprompt "foundry/cli/prompt"
	fpromptCmd "foundry/cli/prompt/cmd"

	"github.com/spf13/cobra"
	"github.com/mlejva/go-prompt"
)

var (
	promptCmd = &cobra.Command{
		Use: 		"prompt",
		Short: 	"",
		Long:		"",
		Run:		runPrompt,
	}

	cmds = []*fprompt.Cmd{
		// fpromptCmd.Watch(),
		fpromptCmd.Exit(),
	}

	saved = false

	col = 0
	row = 0

	outputText = ""
	inputText = ""

	promptw = prompt.NewStandardOutputWriter()
)

func init() {
	rootCmd.AddCommand(promptCmd)
}

func completer(d prompt.Document) []prompt.Suggest {
	col = d.CursorPositionCol()
	inputText = d.CurrentLine()

	s := []prompt.Suggest{}
	for _, c := range cmds {
		s = append(s, c.ToSuggest())
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func executor(s string) {
	if s == "" { return }

	fields := strings.Fields(s)

	if cmd := getCommand(fields[0]); cmd != nil {
		args := fields[1:]

		if err := cmd.Do(args); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {

		promptw.CursorGoTo(row - 2, 0)
		promptw.Flush()
		promptw.WriteRawStr("\x1b[2K") // Erase current line
		promptw.Flush()

		promptw.CursorGoTo(row - 1, 0)
		promptw.Flush()
		t := fmt.Sprintf("Unknown command '%s'. Write 'help' to list available commands.", fields[0])
		promptw.WriteStr(t)
		promptw.Flush()

		// promptw.CursorGoTo(row, 0)
		// promptw.WriteStr("> " + inputText)
		promptw.CursorGoTo(row , col + 3)
		promptw.Flush()
		// fmt.Printf(t)
	}
}

func getCommand(s string) *fprompt.Cmd {
	for _, c := range cmds {
		if c.Text == s {
			return c
		}
	}
	return nil
}

func printPeriodically(ticker *time.Ticker, p *prompt.Prompt) {
	stdoutw := prompt.NewStandardOutputWriter()

	// stdoutw.SaveCursor()
	// stdoutw.Flush()

	for {
		select {
		case <-ticker.C:
			// promptw.SaveCursor()
			// promptw.Flush()

			// stdoutw.Flush()
			if saved {
				stdoutw.UnSaveCursor()
				saved = false
			} else {
				stdoutw.CursorGoTo(0, 0)
			}
			stdoutw.Flush()


			// fmt.Print("\x1b[2k") // Erase current line

			// fmt.Print("\x1b[2J") // Erase screen
			// stdoutw.CursorGoTo(0, 0)
			// stdoutw.Flush()

			// Go to prompt line, erase the line
			// print output text, restore the prompt line
			stdoutw.SaveCursor()
			stdoutw.Flush()
			stdoutw.CursorGoTo(row, 0)
			stdoutw.Flush()
			fmt.Print("\x1b[2K") // Erase current line

			stdoutw.UnSaveCursor()
			stdoutw.Flush()

			// outputText = fmt.Sprintf("%sTick\ncol: %v\n\n", outputText, col)
			t := fmt.Sprintf("Tick\ncol: %v\n", col)
			t += "hello\n"
			t += "\n\n"
			// fmt.Print(outputText + "\n")
			stdoutw.WriteStr(t)
			stdoutw.Flush()

			// promptw.Flush()

			// Save cursor position
			stdoutw.SaveCursor()
			stdoutw.Flush()
			// stdoutw.Flush()
			saved = true

			// stdoutw.Flush()

			// for i := 0; i < 20; i+=1 {
			// 	stdoutw.ScrollDown()
			// 	stdoutw.Flush()
			// }


			promptw.CursorGoTo(row, 0)
			promptw.WriteStr("> " + inputText)
			promptw.CursorGoTo(row, col + 3)
			promptw.Flush()

			// promptw.UnSaveCursor()
			// promptw.Flush()
		}
	}
}

func runPrompt(cmd *cobra.Command, ars []string) {
	interup := prompt.OptionAddKeyBind(prompt.KeyBind{
		Key: 	prompt.ControlC,
		Fn: 	func(buf *prompt.Buffer) {
						os.Exit(0)
					},
	})

	parser := prompt.NewStandardInputParser()
	size := parser.GetWinSize()
	col = int(size.Col)
	row = int(size.Row)


	// w.EraseLine()
	// w.Flush()


	// w.CursorGoTo(int(size.Row) - 3, 0)
	// w.Flush()
	// fmt.Printf("%s", strings.Repeat("-", int(size.Col)))
	// w.WriteStr()

	// getCursorPos(promptw)

	promptw.EraseScreen()
	promptw.CursorGoTo(row, 0)
	promptw.Flush()

	// getCursorPos(promptw)


	p := prompt.New(executor, completer, interup)

	// ticker := time.NewTicker(time.Second * 3)
	// go printPeriodically(ticker, p)

	// d := prompt.NewDocument()

	p.Run()
}


func getCursorPos(cw prompt.ConsoleWriter) (row, col int) {
	fmt.Print("")
	// fmt.Print("\033[6n")

	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// fmt.Print("\x1b[6n")

	// fmt.Fprintf(os.Stdout, "\x1b[6n");
	// fmt.Print("\033[6n")

	cw.AskForCPR()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		fmt.Println(buf)
		fmt.Println(buf.String())
		outC <- buf.String()
	}()

	cw.Flush()

	 // back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC
	// out := buf.String()

	fmt.Println("out:", out)

	fmt.Println("============")
	// s := strings.Split(out, ";")
	// fmt.Println("s:", s)

	// rowStr := s[0][1:]

	// sizeColEl := len(s[1])
	// colStr := s[1][:sizeColEl-1]

	// fmt.Println("rowStr:", rowStr)
	// fmt.Println("colStr:", colStr)

	return 0, 0
}


func toFile() {
	cmd := exec.Command("echo", "-e", "'\x1b[6n'")

	// open the out file for writing
	outfile, err := os.Create("./out.txt")
	if err != nil {
			panic(err)
	}
	defer outfile.Close()
	cmd.Stdout = outfile

	err = cmd.Start(); if err != nil {
			panic(err)
	}
	cmd.Wait()
}

func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		log.Println("buf:", buf)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}

func cursorPos() {
// exec < /dev/tty
// oldstty=$(stty -g)
// stty raw -echo min 0
// echo -en "\033[6n" > /dev/tty
// IFS=';' read -r -d R -a pos
// stty $oldstty
// eval "$1[0]=$((${pos[0]:2} - 2))"
// eval "$1[1]=$((${pos[1]} - 1))"

// eval "$1[0]=$((${pos[0]:2} - 2))"
// eval "$1[1]=$((${pos[1]} - 1))"

	cmd := exec.Command("echo", "-en", "'\x1b[6n'")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("OUTPUT:", string(stdout))
}