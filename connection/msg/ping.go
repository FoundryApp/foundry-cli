package msg

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"foundry/cli/logger"
)

type PingBody struct {
	Token string `json:"token"`
}

type PingMsg struct {
	URL  string
	Body PingBody
}

func NewPingMsg(url, t string) *PingMsg {
	return &PingMsg{
		URL:  url,
		Body: PingBody{t},
	}
}

func (pm *PingMsg) Send() error {
	j, err := json.Marshal(pm.Body)
	if err != nil {
		return err
	}

	res, err := http.Post(pm.URL, "application/json", bytes.NewBuffer(j))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.FdebuglnError("[ping] Error reading ping response body: ", err)
			return err
		}

		bodyString := string(bodyBytes)
		logger.FdebuglnError("[ping] non-ok response:", bodyString)
	}

	return nil
}
