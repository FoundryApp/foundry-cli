package auth

import (
	"strconv"
	"time"

	"foundry/cli/config"
)


func (a *Auth) ClearTokens() {
	a.UserID = ""
	a.Email = ""
	a.IDToken = ""
	a.RefreshToken = ""
	a.ExpiresIn = "0"

	config.Set(idTokenKey, "")
	config.Set(refreshTokenKey, "")
}

func (a *Auth) SaveTokens() error {
	config.Set(idTokenKey, a.IDToken)
	config.Set(refreshTokenKey, a.RefreshToken)

	err := config.Write()

	// if err := os.Setenv(idTokenKey, a.IDToken); err != nil {
	// 	return err
	// }
	// if err := os.Setenv(refreshTokenKey, a.RefreshToken); err != nil {
	// 	return err
	// }
	return err
}

func (a *Auth) LoadTokens() {
	idtok, ok := config.Get(idTokenKey).(string)

	if !ok {
		// TODO: Error
	}
	a.IDToken = idtok

	rtok, ok := config.Get(refreshTokenKey).(string)
	if !ok {
		// TODO: Error
	}
	a.RefreshToken = rtok
}

// Exchanges a refresh token for an ID token
func (a *Auth) RefreshIDToken() error {
	now := time.Now()
	origin := a.originDate

	if a.ExpiresIn == "" { a.ExpiresIn = "0" }

	expireSeconds, err := strconv.Atoi(a.ExpiresIn)
	if err != nil {
		return err
	}

	end := origin.Add(time.Duration(expireSeconds))

	// log.Println("now", now)
	// log.Println("origin", origin)
	// log.Println("expireSeconds", expireSeconds)
	// log.Println("end", end)

	if now.After(end) {
		if err := a.doRefreshReq(); err != nil {
			return err
		}
	}
	return nil
}