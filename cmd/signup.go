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
	signUpCmd = &cobra.Command{
		Use: 		"sign-up",
		Short: 	"Sign up for Foundry in your terminal",
		Long: 	"",
		Run:		runSignUp,
	}

	emailQ = []*survey.Question{
			{
					Name:     "email",
					Prompt:   &survey.Input{Message: "Email:"},
					Validate: survey.Required,
			},
	}

	passQs = []*survey.Question{
		{
			Name: "pass",
			Prompt:   &survey.Password{Message: "Password:"},
			Validate: survey.Required,
		},
		{
			Name: "passAgain",
			Prompt:   &survey.Password{Message: "Password again:"},
			Validate: survey.Required,
		},
	}
)

func init() {
	rootCmd.AddCommand(signUpCmd)
}

func runSignUp(cmd *cobra.Command, args []string) {
	creds := struct {
		Email 		string 	`survey:"email`
		Pass			string	`survey:"pass`
		PassAgain	string	`survey:"passAgain`
	}{}

	// Ask for email
	err := survey.Ask(emailQ, &creds)
	// Without this specific "if" SIGINT (Ctrl+C) would only
	// interrupt the survey's prompt and not the whole program
	if err == terminal.InterruptErr {
		os.Exit(0)
	} else if err != nil {
		log.Println(err)
	}

	// Ask for password
	err = survey.Ask(passQs, &creds)
	// Without this specific "if" SIGINT (Ctrl+C) would only
	// interrupt the survey's prompt and not the whole program
	if err == terminal.InterruptErr {
		os.Exit(0)
	} else if err != nil {
		log.Println(err)
	}

	if creds.Pass != "" && creds.Pass != creds.PassAgain {
		color.Red("\n⨯ Passwords don't match. Please try again.")
		return
	}

	if err = authClient.SignUp(creds.Email, creds.Pass); err != nil {
		color.Red("⨯ Error")
		log.Println("HTTP request error", err)
		return
	}

	if authClient.Error != nil {
		color.Red("⨯ Error")
		log.Println("Auth error", authClient.Error)
		return
	}

	if err = authClient.SaveTokens(); err != nil {
		color.Red("⨯ Error")
		log.Println("Save tokens error", err)
		return
	}

	color.Green("\n✔ Signed up")
}
