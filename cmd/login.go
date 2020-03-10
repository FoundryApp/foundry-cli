package cmd

import (
	// "bytes"
	// "encoding/json"
	// "errors"
	// "fmt"
	// "io/ioutil"
	// "net/http"

	"context"
	"log"
	"os"

	"foundry/cli/auth"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

var (
	loginCmd = &cobra.Command{
		Use: 		"login",
		Short: 	"Log to your Foundry account",
		Long:		"",
		Run:		runLogin,
	}

	qs = []*survey.Question{
			{
					Name:     "email",
					Prompt:   &survey.Input{Message: "Email:"},
					Validate: survey.Required,
			},
			{
					Name: "pass",
					Prompt:   &survey.Input{Message: "Password:"},
					Validate: survey.Required,
			},
	}
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) {
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

	a := auth.New()
	if err = a.SignIn(context.Background(), creds.Email, creds.Pass); err != nil {
		color.Red("⨯ Error")
		log.Println("HTTP request error", err)
		return
	}

	if a.Error != nil {
		color.Red("⨯ Error")
		log.Println("Auth error", a.Error)
		return
	}

	if err = a.SaveTokens(); err != nil {
		color.Red("⨯ Error")
		log.Println("Save tokens error", err)
		return
	}

	// if err = runHTTPPost(a.IDToken); err != nil {
	// 	color.Red("⨯ Error")
	// 	log.Println("HTTP post error", err)
	// 	return
	// }

	color.Green("✔ Signed in")
}

// func runHTTPPost(idToken string) error {
// 	req := struct {
// 		Token string `json:"token"`
// 	}{idToken}

// 	jReq, err := json.Marshal(req)
// 	if err != nil {
// 		return err
// 	}

// 	url := fmt.Sprintf("http://127.0.0.1:8081/run")
// 	res, err := http.Post(url, "application/json", bytes.NewBuffer(jReq))
// 	if err != nil {
// 		return err
// 	}

// 	if res.StatusCode != http.StatusOK {
// 		body, _ := ioutil.ReadAll(res.Body)
// 		return errors.New(string(body))
// 	}

// 	return nil
// }
