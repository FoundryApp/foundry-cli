package msg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"foundry/cli/connection/endpoint"
	"foundry/cli/logger"
)

type Env struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type EnvBody struct {
	Token string `json:"token"`
	Envs  []Env  `json:"envs"`
}

type EnvMsg struct {
	url  string
	Body EnvBody
}

func SetEnvURL() string {
	return fmt.Sprintf("%s://%s/setenv", endpoint.SetEnvScheme, endpoint.SetEnvURL)
}

func NewEnvMsg(token string, envs []Env) *EnvMsg {
	return &EnvMsg{
		url:  SetEnvURL(),
		Body: EnvBody{token, envs},
	}
}

func (em *EnvMsg) Send() error {
	j, err := json.Marshal(em.Body)
	if err != nil {
		return err
	}

	res, err := http.Post(em.url, "application/json", bytes.NewBuffer(j))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.FdebuglnError("Error reading env msg response body: ", err)
			return err
		}

		bodyString := string(bodyBytes)
		logger.FdebuglnError("Non-ok env msg response:", bodyString)
	}
	return nil
}
