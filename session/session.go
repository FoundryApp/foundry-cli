package session

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
)

type Session struct {
	WebSocketHost string `json:"websocketCli"`
	wsconn        *websocket.Conn
}

func GetCurrent() (*Session, error) {
	url := "http://localhost:3600/session"
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	decoded := struct {
		Sess *Session `json:"session"`
		Err  string   `json:"error"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		return nil, err
	}

	if decoded.Err != "" {
		return nil, errors.New(decoded.Err)
	}
	return decoded.Sess, nil
}
