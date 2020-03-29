package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/fatih/color"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

var (
	signInCmd = &cobra.Command{
		Use: 		"sign-in",
		Short: 	"Sign in to your Foundry account",
		Long:		"",
		Run:		runSignIn,
	}

	qs = []*survey.Question{
			{
					Name:     "email",
					Prompt:   &survey.Input{Message: "Email:"},
					Validate: survey.Required,
			},
			{
					Name: 		"pass",
					Prompt:   &survey.Password{Message: "Password:"},
					Validate: survey.Required,
			},
	}
)

func init() {
	rootCmd.AddCommand(signInCmd)
}

func runSignIn(cmd *cobra.Command, args []string) {
	creds := struct {
		Email string 	`survey:"email`
		Pass	string	`survey:"pass`
	}{}

	err := survey.Ask(qs, &creds)
	// Without this specific "if" SIGINT (Ctrl+C) would only
	// interrupt the survey's prompt and not the whole program
	if err == terminal.InterruptErr {
		os.Exit(0)
	} else if err != nil {
		log.Println(err)
	}

	if err = authClient.SignIn(creds.Email, creds.Pass); err != nil {
		color.Red("⨯ Error")
		log.Println("HTTP request error", err)
		return
	}

	if authClient.Error != nil {
		color.Red("⨯ Error")
		log.Println("Auth error", authClient.Error)
		return
	}

	color.Green("✔ Signed in")
}
