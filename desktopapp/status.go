package desktopapp

import (
	"encoding/json"
	"net/http"
)

type AppStatus string

const (
	ReadyAppStatus    AppStatus = "ready"
	NotReadyAppStatus AppStatus = "not-ready"
)

func GetStatus() (AppStatus, error) {
	url := "http://localhost:3600/status"
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	decoded := struct {
		Status AppStatus `json:"status"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		return "", err
	}

	return decoded.Status, nil
}
