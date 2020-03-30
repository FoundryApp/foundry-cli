package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"foundry/cli/logger"
)

const (
	apiKey 					= "AIzaSyAqL--IsyZd3cQTUgXR3KRWZZN-M6jR1kE"
	idTokenKey			= "FOUNDRY_AUTH_ID_TOKEN"
	refreshTokenKey = "FOUNDRY_AUTH_REFRESH_TOKEN"
	authStateKey 		= "FOUNDRY_AUTH_STATE"
)

type AuthStateType int

// WARNING: It's important the order doesn't change because the AuthState field on Auth struct
// is serialized in the config file.
// Changing the order of the following consts would cause that serialized values would have
// a different logical meaning
const (
	AuthStateTypeSignedOut					AuthStateType = iota + 1 // +1 so the first const's value is different from zero int value
	AuthStateTypeSignedIn
	AuthStateTypeSignedInAnonymous
)

type AuthError struct {
	Message 		string	`json:"message"`
	StatusCode 	int			`json:"code"`
}

func (ae *AuthError) Error() string {
	return fmt.Sprintf("[%v] %v\n", ae.StatusCode, ae.Message)
}

// TODO: Find out if we can serialize structures using viper
type Auth struct {
	Error					*AuthError	`json:"error"`
	UserID				string			`json:"localId"`
	Email					string			`json:"email"`
	IDToken				string			`json:"idToken"`
	RefreshToken 	string 			`json:"refreshToken"`

	ExpiresIn			string			`json:"expiresIn"`
	originDate		time.Time

	AuthState			AuthStateType

	DisplayName		string			`json:"displayName"`
}

func New() (*Auth, error) {
	a := &Auth{
		AuthState: AuthStateTypeSignedOut,
	}
	if err := a.loadTokensAndState(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Auth) SignUp(email, pass string) error {
	baseURL := "https://identitytoolkit.googleapis.com/v1"
	var endpoint string
	var reqBody interface{}

	if a.AuthState == AuthStateTypeSignedInAnonymous {
		// Check if auth state is AuthStateTypeSignedInAnonymous
		// If so, link the anonymous user with email, password, and IDToken
		logger.Fdebugln("Signing up an anonymous user (= linking email + pass)")
		endpoint = fmt.Sprintf("accounts:update?key=%v", apiKey)
		reqBody = struct{
			IDToken						string 	`json:"idToken"`
			Email							string	`json:"email"`
			Password					string	`json:"password"`
			ReturnSecureToken	bool		`json:"returnSecureToken"`
		}{a.IDToken, email, pass, true}
	} else {
		logger.Fdebugln("Signing up a new user")
		endpoint = fmt.Sprintf("accounts:signUp?key=%v", apiKey)
		reqBody = struct{
			Email							string	`json:"email"`
			Password					string	`json:"password"`
			ReturnSecureToken	bool		`json:"returnSecureToken"`
		}{email, pass, true}
	}

	url := fmt.Sprintf("%v/%v", baseURL, endpoint)

	if err := a.doAuthReq(url, reqBody); err != nil {
		return err
	}

	if a.Error != nil {
		return nil
	}

	oldState := a.AuthState
	a.AuthState = AuthStateTypeSignedIn
	if err := a.saveTokensAndState(); err != nil {
		a.AuthState = oldState
		return err
	}
	return nil
}

func (a *Auth) SignUpAnonymously() error {
	reqBody := struct {
		ReturnSecureToken	bool `json:"returnSecureToken"`
	}{true}

	baseURL := "https://identitytoolkit.googleapis.com/v1"
	endpoint := fmt.Sprintf("accounts:signUp?key=%v", apiKey)
	url := fmt.Sprintf("%v/%v", baseURL, endpoint)

	if err := a.doAuthReq(url, reqBody); err != nil {
		return err
	}

	if a.Error != nil {
		return nil
	}

	oldState := a.AuthState
	a.AuthState = AuthStateTypeSignedInAnonymous
	if err := a.saveTokensAndState(); err != nil {
		a.AuthState = oldState
		return err
	}
	return nil
}

func (a *Auth) SignIn(email, pass string) error {
	reqBody := struct {
		Email 						string 	`json:"email"`
		Password 					string 	`json:"password"`
		ReturnSecureToken bool 		`json:"returnSecureToken"`
	}{email, pass, true}

	baseURL := "https://identitytoolkit.googleapis.com/v1"
	endpoint := fmt.Sprintf("accounts:signInWithPassword?key=%v", apiKey)
	url := fmt.Sprintf("%v/%v", baseURL, endpoint)

	if err := a.doAuthReq(url, reqBody); err != nil {
		return err
	}

	if a.Error != nil {
		return nil
	}

	oldState := a.AuthState
	a.AuthState = AuthStateTypeSignedIn
	if err := a.saveTokensAndState(); err != nil {
		a.AuthState = oldState
		return err
	}

	return nil
}

func (a *Auth) SignOut() error {
	a.Error = nil
	a.UserID = ""
	a.Email = ""
	a.IDToken = ""
	a.RefreshToken = ""
	a.ExpiresIn = "0"
	a.AuthState = AuthStateTypeSignedOut
	return a.clearTokensAndState()
}

func (a *Auth) doAuthReq(url string, body interface{}) error {
	a.Error = nil

	jBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jBody))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, a)
	if err != nil {
		return err
	}

	if a.Error != nil {
		return nil
	}

	// Save the time when we originaly acquired the ID token
	// for checking whether we need to refresh it
	a.originDate = time.Now()

	return nil
}

func (a *Auth) doRefreshReq() error {
	logger.Fdebugln("Refreshing ID token")

	u := fmt.Sprintf("https://securetoken.googleapis.com/v1/token?key=%v", apiKey)
	data := url.Values{}
	data.Set("refresh_token", a.RefreshToken)
	data.Set("grant_type", "refresh_token")

	req, err := http.NewRequest("POST", u, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Sigh... Firebase has different keys in the response payload
	// for token refresh flow from the response payload in a sign
	// in flow. Also, its content-type isn't application/json but
	// application/x-www-form-urlencoded.
	var j struct {
		ExpiresIn 		string `json:"expires_in"`
		RefreshToken 	string `json:"refresh_token"`
		IDToken				string `json:"id_token"`
	}

	err = json.Unmarshal(bodyBytes, &j)
	if err != nil {
		return err
	}

	a.ExpiresIn = j.ExpiresIn
	a.IDToken = j.IDToken
	a.RefreshToken = j.RefreshToken
	a.originDate = time.Now()

	// TODO: error checking
	// TOKEN_EXPIRED: The user's credential is no longer valid. The user must sign in again.
	// USER_DISABLED: The user account has been disabled by an administrator.
	// USER_NOT_FOUND: The user corresponding to the refresh token was not found. It is likely the user was deleted.
	// API key not valid. Please pass a valid API key. (invalid API key provided)
	// INVALID_REFRESH_TOKEN: An invalid refresh token is provided.
	// Invalid JSON payload received. Unknown name \"refresh_tokens\": Cannot bind query parameter. Field 'refresh_tokens' could not be found in request message.
	// INVALID_GRANT_TYPE: the grant type specified is invalid.
	// MISSING_REFRESH_TOKEN: no refresh token provided.

	return nil
}