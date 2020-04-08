package firebase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Error struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type Request struct {
	Data interface{} `json:"data"`
}

type Response struct {
	Error  *Error      `json:"error"`
	Result interface{} `json:"result"`
}

func Call(funcName, IDToken string, data interface{}) (*Response, error) {
	url := fmt.Sprintf("https://us-central1-foundryapp.cloudfunctions.net/%s", funcName)

	var reqBody Request
	if data == nil {
		// Firebase httpsCallable functions requires that there's always at least
		// empty 'data' field (e.i.: '"data": {}') in the body
		reqBody = Request{struct{}{}}
	} else {
		reqBody = struct {
			Data interface{} `json:"data"`
		}{data}
	}

	marshaledBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshaledBody))
	if err != nil {
		return nil, err
	}

	bearer := "Bearer " + IDToken
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 30}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var respBody Response
	if err = json.Unmarshal(bodyBytes, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}
