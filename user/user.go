package user

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Credentials struct {
	Token        string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

type User struct {
	ID    string      `json:"uid"`
	Creds Credentials `json:"creds"`
}

func GetCurrent() (*User, error) {
	url := "http://localhost:3600/user"
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoded := struct {
		User *User  `json:"user"`
		Err  string `json:"error"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		return nil, err
	}

	if decoded.Err != "" {
		return nil, errors.New(decoded.Err)
	}
	return decoded.User, nil
}
