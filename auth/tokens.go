package auth

import (
	"fmt"
	"strconv"
	"time"

	"foundry/cli/config"
)

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

func (a *Auth) saveTokens() error {
	config.Set(idTokenKey, a.IDToken)
	config.Set(refreshTokenKey, a.RefreshToken)
	return config.Write()
}

func (a *Auth) loadTokens() error {
	idtok, ok := config.Get(idTokenKey).(string)

	if !ok {
		return fmt.Errorf("Failed to get ID token from config")
	}
	a.IDToken = idtok

	rtok, ok := config.Get(refreshTokenKey).(string)
	if !ok {
		return fmt.Errorf("Failed to get refresh token from config")
	}
	a.RefreshToken = rtok

	return nil
}

func (a *Auth) clearTokens() error {
	config.Set(idTokenKey, "")
	config.Set(refreshTokenKey, "")
	return config.Write()
}
