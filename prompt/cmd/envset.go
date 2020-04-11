package cmd

import (
	"fmt"
	c "foundry/cli/connection"
	"foundry/cli/connection/msg"
	"foundry/cli/logger"
	"strings"

	goprompt "github.com/mlejva/go-prompt"
)

type EnvSetCmd struct {
	Text    string
	Desc    string
	RunCh   RunChannelType
	IDToken string
}

func NewEnvSetCmd(IDToken string) *EnvSetCmd {
	return &EnvSetCmd{
		Text:    "env-set",
		Desc:    "Set environment variable(s) in your cloud environment",
		RunCh:   make(chan Args),
		IDToken: IDToken,
	}
}

// Implement Cmd interface
func (c *EnvSetCmd) Run(conn *c.Connection, args Args) (promptOutput string, promptInfo string, err error) {
	if len(args) == 0 {
		return "", "No envs specified. Example usage: 'foundry env-set MY_ENV=ENV_VALUE ANOTHER_ENV=ANOTHER_VALUE'", nil
	}

	envs := []msg.Env{}
	for _, env := range args {
		arr := strings.Split(env, "=")

		if len(arr) != 2 {
			logger.FdebuglnFatal("Error parsing environment variable:", env)
			return "", "", fmt.Errorf(fmt.Sprintf("error parsing environment variable. Expected format 'env=value'. Got: %s", env))
		}

		name := arr[0]
		val := arr[1]

		if name == "" {
			logger.FdebuglnFatal("Error parsing environment variable - name is empty:", env)
			return "", "", fmt.Errorf(fmt.Sprintf("error parsing environment variable. Expected format 'env=value'. Got: %s", env))
		}
		if val == "" {
			logger.FdebuglnFatal("Error parsing environment variable - val is empty:", env)
			return "", "", fmt.Errorf(fmt.Sprintf("error parsing environment variable. Expected format 'env=value'. Got: %s", env))
		}

		envs = append(envs, msg.Env{name, val})
	}

	envMsg := msg.NewEnvMsg(c.IDToken, envs)
	if err := envMsg.Send(); err != nil {
		logger.FdebuglnError("Error setting environment variables:", err)
		return "", "", err
	}
	return "", "Variables set", nil
}

func (c *EnvSetCmd) RunRequest(args Args) {
	c.RunCh <- args
}

func (c *EnvSetCmd) ToSuggest() goprompt.Suggest {
	return goprompt.Suggest{c.Text, c.Desc}
}

func (c *EnvSetCmd) Name() string {
	return c.Text
}

func (c *EnvSetCmd) String() string {
	return fmt.Sprintf("%s - %s", c.Text, c.Desc)
}
