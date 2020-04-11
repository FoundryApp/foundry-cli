package cmd

import (
	"fmt"
	c "foundry/cli/connection"
	"foundry/cli/connection/msg"
	"foundry/cli/firebase"
	"foundry/cli/logger"
	"strings"

	goprompt "github.com/mlejva/go-prompt"
)

type EnvDelCmd struct {
	Text    string
	Desc    string
	RunCh   RunChannelType
	IDToken string
}

func NewEnvDelCmd(IDToken string) *EnvDelCmd {
	return &EnvDelCmd{
		Text:    "env-delete",
		Desc:    "Delete environment variable(s) from your cloud environment",
		RunCh:   make(chan Args),
		IDToken: IDToken,
	}
}

// Implement Cmd interface
func (c *EnvDelCmd) Run(conn *c.Connection, args Args) (promptOutput string, promptInfo string, err error) {
	if len(args) == 0 {
		return "", "No envs to delete specified. Example usage: 'foundry env-delete ENV_1 ENV_2'", nil
	}

	reqBody := struct {
		Delete []string `json:"delete"`
	}{args}
	res, err := firebase.Call("deleteUserEnvs", c.IDToken, reqBody)
	if err != nil {
		logger.FdebuglnFatal("Error calling deleteUserEnvs:", err)
		return "", "", fmt.Errorf(fmt.Sprintf("error deleting environment variables: %s", err))
	}
	if res.Error != nil {
		logger.FdebuglnFatal("Error calling deleteUserEnvs:", res.Error)
		return "", "", fmt.Errorf(fmt.Sprintf("error deleting environment variables: %s", err))
	}

	// Report new envs to Autorun
	envsMap, ok := res.Result.(map[string]interface{})
	if !ok {
		logger.FdebuglnFatal("Failed to type assert res.Result")
		return "", "", fmt.Errorf("error deleting environment variables")
	}

	envs := []msg.Env{}
	for name, val := range envsMap {
		envs = append(envs, msg.Env{name, val.(string)})
	}
	logger.Fdebugln("Sending new envs vars to Autorun:", envs)

	envMsg := msg.NewEnvMsg(c.IDToken, envs)
	if err = envMsg.Send(); err != nil {
		logger.FdebuglnError("Failed to report new env vars (after deletion) to Autorun", err)
	}

	return "", "Deleted " + strings.Join(args, ", "), err
}

func (c *EnvDelCmd) RunRequest(args Args) {
	c.RunCh <- args
}

func (c *EnvDelCmd) ToSuggest() goprompt.Suggest {
	return goprompt.Suggest{c.Text, c.Desc}
}

func (c *EnvDelCmd) Name() string {
	return c.Text
}

func (c *EnvDelCmd) String() string {
	return fmt.Sprintf("%s - %s", c.Text, c.Desc)
}
