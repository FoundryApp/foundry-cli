package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

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
	promptRow = 0

	outputText = ""
	inputText = ""
	errorText = ""

	filledLines = 0

	freeLines = 0
	overlapping = false

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

	return []prompt.Suggest{}
	// return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
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

		promptw.CursorGoTo(promptRow - 1, 0)
		promptw.Flush()
		promptw.EraseLine()
		promptw.Flush()

		errorText = fmt.Sprintf("Unknown command '%s'. Write 'help' to list available commands.", fields[0])
		fmt.Println(errorText)

		promptw.CursorGoTo(promptRow, 0)
		promptw.Flush()

		// promptw.WriteRawStr("\x1b[2K") // Erase current line
		// promptw.Flush()

		// promptw.CursorGoTo(row - 1, 0)
		// // promptw.Flush()
		// t := fmt.Sprintf("Unknown command '%s'. Write 'help' to list available commands.", fields[0])
		// promptw.WriteStr(t)
		// // promptw.Flush()

		// // promptw.CursorGoTo(row, 0)
		// // promptw.WriteStr("> " + inputText)
		// promptw.CursorGoTo(row - 1 , col + 3)
		// promptw.Flush()
		// // fmt.Printf(t)
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

func printPeriodically2(ticker *time.Ticker, p *prompt.Prompt) {
	stdoutw := prompt.NewStandardOutputWriter()

	i := 0

	for {
		select {
		case <-ticker.C:
			i += 1
			if saved {
				stdoutw.UnSaveCursor()
				if overlapping {
					stdoutw.CursorUp(2)
				}
				saved = false
			} else {
				stdoutw.CursorGoTo(0, 0)
			}
			stdoutw.Flush()

			stdoutw.SaveCursor()
			stdoutw.Flush()

			// Go to prompt line, erase the line
			// print output text, restore the prompt line
			stdoutw.CursorGoTo(promptRow, 0)
			stdoutw.Flush()
			stdoutw.EraseLine()
			stdoutw.Flush()

			// same with error line
			// Go to error line, erase the line
			if len(errorText) > 0 {
				stdoutw.CursorGoTo(promptRow - 1, 0)
				stdoutw.Flush()
				stdoutw.EraseLine()
				stdoutw.Flush()
			}

			// Restore cursor
			stdoutw.UnSaveCursor()
			stdoutw.Flush()

			// Output the text
			t := fmt.Sprintf("=Lorem \nipsum \ndolor \nsit \namet\n. Hello \n WOrld\n, how - %v\n", i)
			calcOverlap(t)
			stdoutw.WriteStr(t)
			stdoutw.Flush()
			stdoutw.SaveCursor()
			saved = true
			stdoutw.Flush()

			// DO THE FOLLOWING ONLY IF THE OUTPUT TEXT IS ABOUT TO HIT ERROR + PROMPT LINE
			// vars:
			// visibleRows
			// filledRows (# of filled rows of the visible terminal)
			// leftRows (# of rows left, until we hit error or prompt line)

			// OR IDEA: all I need to do is to move cursor 2 rows up next time I'll be outputing the text?
			// I tested this and it will work but I need to do this only when the 2 new lines have pushed
			// the terminal down

			// Create space for prompt line + error line
			stdoutw.WriteRawStr("\n\n")
			stdoutw.Flush()

			// Restore the error line
			promptw.CursorGoTo(promptRow - 1, 0)
			promptw.Flush()
			promptw.WriteStr(errorText)
			promptw.Flush()

			// Restore the prompt line
			promptw.CursorGoTo(promptRow, 0)
			promptw.Flush()
			promptw.WriteStr("> " + inputText)
			promptw.Flush()
		}
	}
}

func calcOverlap(t string) {
	// TODO: Handle case when len(t) > width of terminal - t gets rendered as multiple lines
	l := strings.Split(t, "\n")

	if len(l) >= freeLines {
		freeLines = 0
		overlapping = true
	} else {
		freeLines -= len(l)
	}
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



			// t := "Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nInteger nec odio. Praesent libero.\nSed cursus ante dapibus diam. Sed nisi.\nNulla quis sem at nibh.\nelementum imperdiet.\nDuis sagittis ipsum. Praesent mauris.\nFusce\n"
			t := "Lorem \nipsum \ndolor \nsit \namet\n"
			// t := fmt.Sprintf("Tick\ncol: %v\n", col)
			// if len(t) % 2 == 0 {
			// 	t += "\n"
			// }

			// outputText += t
			// lines := strings.Split(outputText, "\n")

			lines := strings.Split(t, "\n")
			filledLines += len(lines)


			if filledLines > promptRow - 3 {
				stdoutw.SaveCursor()
				stdoutw.Flush()

				// Go to prompt line, erase the line
				// print output text, restore the prompt line
				stdoutw.CursorGoTo(promptRow, 0)
				stdoutw.Flush()
				fmt.Print("\x1b[2K") // Erase current line


				// same with error line
				// Go to error line, erase the line
				if len(errorText) > 0 {
					stdoutw.CursorGoTo(promptRow - 1, 0)
					stdoutw.Flush()
					fmt.Print("\x1b[2K") // Erase current line
				}

				stdoutw.UnSaveCursor()
				stdoutw.Flush()

				// if filledLines > promptRow - 3 {
				// 	fmt.Println("filledLines:", filledLines)
				// 	fmt.Println("(filledLines - promptRow - 3):", (filledLines - promptRow - 3))
				// 	for i := 0; i <= (filledLines); i++ {
				// 		filledLines -= 1
				// 		stdoutw.ScrollDown()
				// 		stdoutw.Flush()
				// 	}
				// 	fmt.Println("filledLines:", filledLines)
				// 	stdoutw.SaveCursor()
				// 	stdoutw.Flush()
				// }
			}
			stdoutw.WriteStr(t)
			stdoutw.Flush()
			stdoutw.SaveCursor()
			saved = true
			stdoutw.Flush()




			// Restore the error line
			promptw.CursorGoTo(promptRow - 1, 0)
			promptw.Flush()
			promptw.WriteStr(errorText)
			promptw.Flush()


			// Restore the prompt line
			promptw.CursorGoTo(promptRow, 0)
			promptw.Flush()
			promptw.WriteStr("> " + inputText)
			// promptw.CursorGoTo(promptRow, col + 3)
			promptw.Flush()




			// fmt.Print("\x1b[2k") // Erase current line

			// fmt.Print("\x1b[2J") // Erase screen
			// stdoutw.CursorGoTo(0, 0)
			// stdoutw.Flush()

			// Go to prompt line, erase the line
			// print output text, restore the prompt line
			// stdoutw.SaveCursor()
			// stdoutw.Flush()
			// stdoutw.CursorGoTo(row, 0)
			// stdoutw.Flush()
			// fmt.Print("\x1b[2K") // Erase current line

			// stdoutw.UnSaveCursor()
			// stdoutw.Flush()

			// outputText = fmt.Sprintf("%sTick\ncol: %v\n\n", outputText, col)
			// t := fmt.Sprintf("Tick\ncol: %v\n", col)
			// t += "hello\n"
			// t += "\n\n\n"
			// // fmt.Print(outputText + "\n")
			// stdoutw.WriteStr(t)
			// stdoutw.Flush()

			// promptw.Flush()

			// Save cursor position
			// stdoutw.SaveCursor()
			// stdoutw.Flush()
			// stdoutw.Flush()
			// saved = true

			// stdoutw.Flush()

			// stdoutw.ScrollDown()
			// stdoutw.Flush()

			// promptw.CursorGoTo(row, 0)
			// promptw.Flush()
			// promptw.WriteStr("> " + inputText)
			// promptw.Flush()
			// promptw.CursorGoTo(row, col + 3)
			// promptw.Flush()

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

	promptRow = row
	freeLines = promptRow - 3


	promptw.EraseScreen()
	promptw.CursorGoTo(promptRow, 0)
	promptw.Flush()

	p := prompt.New(executor, completer, interup)


	ticker := time.NewTicker(time.Second * 1)
	go printPeriodically2(ticker, p)

	p.Run()
}
