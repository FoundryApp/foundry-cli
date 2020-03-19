package cmd

import (
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

var (
	promptCmd = &cobra.Command{
		Use: 		"prompt",
		Short: 	"",
		Long:		"",
		Run:		runPrompt,
	}
)

func init() {
	rootCmd.AddCommand(promptCmd)
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "watch", Description: "Watch specific function(s)"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func executor(s string) {
	// fmt.Println(s)
}

func runPrompt(cmd *cobra.Command, ars []string) {
	fmt.Println("Running prompt2")
	interup := prompt.OptionAddKeyBind(prompt.KeyBind{
		Key: 	prompt.ControlC,
		Fn: 	func(buf *prompt.Buffer) {
						os.Exit(0)
					},
	})
	// prompt.Input("> ", completer, interr)
	p := prompt.New(executor, completer, interup)
	p.Run()
	// for {

		// t := prompt.Input("> ", completer, interr)
		// fmt.Printf("%s\n", t)
	// }
}