package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	apiKey 					= "AIzaSyAqL--IsyZd3cQTUgXR3KRWZZN-M6jR1kE"
	baseURL 				= "https://identitytoolkit.googleapis.com/v1"
	idTokenKey			= "FOUNDRY_AUTH_ID_TOKEN"
	refreshTokenKey = "FOUNDRY_AUTH_REFRESH_TOKEN"
)

type AuthError struct {
	Message 		string	`json:"message"`
	StatusCode 	int			`json:"code"`
}

func (ae *AuthError) Error() string {
	return fmt.Sprintf("[%v] %v\n", ae.StatusCode, ae.Message)
}

type Auth struct {
	Error					*AuthError	`json:"error"`
	UserID				string			`json:"localId"`
	Email					string			`json:"email"`
	IDToken				string			`json:"idToken"`
	RefreshToken 	string 			`json:"refreshToken"`
}

func New() *Auth {
	return &Auth{}
}

func (a *Auth) SignIn(ctx context.Context, email, pass string) error {
	url := fmt.Sprintf("%v/accounts:signInWithPassword?key=%v", baseURL, apiKey)

	req := struct {
		Email 						string 	`json:"email"`
		Password 					string 	`json:"password"`
		ReturnSecureToken bool 		`json:returnSecureToken`
	}{email, pass, true}
	jReq, err := json.Marshal(req)
	if err != nil {
		return err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jReq))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, a)
	return err
}

func (a *Auth) SaveTokens() error {
	if err := os.Setenv(idTokenKey, a.IDToken); err != nil {
		return err
	}

	if err := os.Setenv(refreshTokenKey, a.RefreshToken); err != nil {
		return err
	}
	return nil
}

func (a *Auth) LoadTokens() {
	a.IDToken = os.Getenv(idTokenKey)
	a.RefreshToken = os.Getenv(refreshTokenKey)
}