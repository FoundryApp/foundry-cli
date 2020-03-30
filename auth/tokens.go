package auth

import (
	"strconv"
	"time"

	"foundry/cli/config"
	"foundry/cli/logger"
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
		if err := a.saveTokensAndState(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Auth) saveTokensAndState() error {
	config.Set(idTokenKey, a.IDToken)
	config.Set(refreshTokenKey, a.RefreshToken)
	config.Set(authStateKey, a.AuthState)
	return config.Write()
}

func (a *Auth) loadTokensAndState() error {
	idtok := config.GetString(idTokenKey)
	a.IDToken = idtok

	rtok := config.GetString(refreshTokenKey)
	a.RefreshToken = rtok

	state := config.GetInt(authStateKey)
	// State is 0 when the config file is empty
	if state != 0 {
		a.AuthState = AuthStateType(state)
	} else {
		a.AuthState = AuthStateTypeSignedOut
	}

	logger.Fdebugln("Loaded AuthState from config (1 = signed out, 2 = signed in, 3 = anonymous):", a.AuthState)

	return nil
}

func (a *Auth) clearTokensAndState() error {
	config.Set(idTokenKey, "")
	config.Set(refreshTokenKey, "")
	config.Set(authStateKey, AuthStateTypeSignedOut)
	return config.Write()
}
