package cmd

import (
	"fmt"
	"os"
	"strings"

	fprompt "foundry/cli/prompt"
	fpromptCmd "foundry/cli/prompt/cmd"

	"github.com/spf13/cobra"
	"github.com/c-bata/go-prompt"
)

var (
	promptCmd = &cobra.Command{
		Use: 		"prompt",
		Short: 	"",
		Long:		"",
		Run:		runPrompt,
	}

	cmds = []*fprompt.Cmd{
		fpromptCmd.Watch(),
		fpromptCmd.Exit(),
	}
)

func init() {
	rootCmd.AddCommand(promptCmd)
}

func completer(d prompt.Document) []prompt.Suggest {
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
		fmt.Printf("Unknown command '%s'. Write 'help' to list available commands.\n", fields[0])
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

func runPrompt(cmd *cobra.Command, ars []string) {
	interup := prompt.OptionAddKeyBind(prompt.KeyBind{
		Key: 	prompt.ControlC,
		Fn: 	func(buf *prompt.Buffer) {
						os.Exit(0)
					},
	})

	p := prompt.New(executor, completer, interup)
	p.Run()
}