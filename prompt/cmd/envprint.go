package cmd

import (
	"fmt"
	c "foundry/cli/connection"
	"foundry/cli/firebase"
	"foundry/cli/logger"

	goprompt "github.com/mlejva/go-prompt"
)

type EnvPrintCmd struct {
	Text    string
	Desc    string
	RunCh   RunChannelType
	IDToken string
}

func NewEnvPrintCmd(IDToken string) *EnvPrintCmd {
	return &EnvPrintCmd{
		Text:    "env-print",
		Desc:    "Print all environment variables in your cloud environment",
		RunCh:   make(chan Args),
		IDToken: IDToken,
	}
}

// Implement Cmd interface
func (c *EnvPrintCmd) Run(conn *c.Connection, args Args) (promptOutput string, promptInfo string, err error) {
	res, err := firebase.Call("getUserEnvs", c.IDToken, nil)
	if err != nil {
		logger.FdebuglnFatal("Error calling getUserEnvs:", err)
		return "", "", err
	}
	if res.Error != nil {
		logger.FdebuglnFatal("Error calling getUserEnvs:", res.Error)
		return "", "", fmt.Errorf(res.Error.Message)
	}

	envs, ok := res.Result.(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("error printing environment variables. Failed to convert the response")
	}

	if len(envs) == 0 {
		return "", "No environment variable has been set yet", nil
	}

	delimiter := "-----------------------------------------------------"
	msg := "\n" + delimiter + "\n|\n"
	msg += "| Following environment variables are set:\n|"
	for k, v := range envs {
		msg += fmt.Sprintf("\n|  %s=%s", k, v.(string))
	}
	msg += "\n|\n" + delimiter + "\n"

	return msg, "", nil
}

func (c *EnvPrintCmd) RunRequest(args Args) {
	c.RunCh <- args
}

func (c *EnvPrintCmd) ToSuggest() goprompt.Suggest {
	return goprompt.Suggest{c.Text, c.Desc}
}

func (c *EnvPrintCmd) Name() string {
	return c.Text
}

func (c *EnvPrintCmd) String() string {
	return fmt.Sprintf("%s - %s", c.Text, c.Desc)
}
