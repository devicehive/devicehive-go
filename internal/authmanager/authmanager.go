package authmanager

import (
	"encoding/json"
	"github.com/devicehive/devicehive-go/internal/requester"
	"github.com/devicehive/devicehive-go/internal/resourcenames"
)

func New(reqstr requester.Requester) *AuthManager {
	return &AuthManager{
		Reauth: &ReauthenticationState{},
		reqstr: reqstr,
	}
}

type AuthManager struct {
	Reauth       *ReauthenticationState
	reqstr       requester.Requester
	accessToken  string
	login        string
	password     string
	refreshToken string
}

func (a *AuthManager) SetCreds(login, password string) {
	a.login = login
	a.password = password
}

func (a *AuthManager) SetRefreshToken(refTok string) {
	a.refreshToken = refTok
}

func (a *AuthManager) SetAccessToken(accTok string) {
	a.accessToken = accTok
}

func (a *AuthManager) AccessToken() string {
	return a.accessToken
}

func (a *AuthManager) RefreshToken() (accessToken string, err error) {
	if a.refreshToken == "" {
		accessToken, _, err = a.tokensByCreds(a.login, a.password)
		return accessToken, err
	}

	return a.accessTokenByRefresh(a.refreshToken)
}

func (a *AuthManager) tokensByCreds(login, pass string) (accessToken, refreshToken string, err error) {
	rawRes, err := a.request(resourcenames.TokenByCreds, map[string]interface{}{
		"login":    login,
		"password": pass,
	})

	if err != nil {
		return "", "", err
	}

	tok := &token{}
	parseErr := json.Unmarshal(rawRes, tok)

	if parseErr != nil {
		return "", "", parseErr
	}

	return tok.Access, tok.Refresh, nil
}

func (a *AuthManager) accessTokenByRefresh(refreshToken string) (accessToken string, err error) {
	rawRes, err := a.request(resourcenames.TokenRefresh, map[string]interface{}{
		"refreshToken": refreshToken,
	})

	if err != nil {
		return "", err
	}

	tok := &token{}
	parseErr := json.Unmarshal(rawRes, tok)

	if parseErr != nil {
		return "", parseErr
	}

	return tok.Access, nil
}

func (a *AuthManager) request(resourceName string, data map[string]interface{}) ([]byte, error) {
	return a.reqstr.Request(resourceName, data, 0, a.accessToken)
}
