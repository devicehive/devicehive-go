package dh

import (
	"encoding/json"
	"time"
)

type token struct {
	Access  string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}

func (c *Client) CreateToken(userId int, expiration time.Time, actions, networkIds, deviceTypeIds []string) (accessToken, refreshToken string, err *Error) {
	payload := map[string]interface{}{
		"userId": userId,
	}

	if actions != nil {
		payload["actions"] = actions
	}
	if networkIds != nil {
		payload["networkIds"] = networkIds
	}
	if deviceTypeIds != nil {
		payload["deviceTypeIds"] = deviceTypeIds
	}
	if expiration.Unix() > 0 {
		payload["expiration"] = expiration.UTC().Format(timestampLayout)
	}

	data := map[string]interface{}{
		"action":  "token/create",
		"payload": payload,
	}

	return c.tokenRequest(data)
}

func (c *Client) RefreshToken() (accessToken string, err *Error) {
	if c.refreshToken == "" {
		accessToken, _, err = c.tokensByCreds(c.login, c.password)

		if err != nil {
			return "", err
		}

		return accessToken, nil
	}

	return c.accessTokenByRefresh(c.refreshToken)
}

func (c *Client) accessTokenByRefresh(refreshToken string) (accessToken string, err *Error) {
	_, resBytes, err := c.request(map[string]interface{}{
		"action":       "token/refresh",
		"refreshToken": c.refreshToken,
	})

	if err != nil {
		return "", err
	}

	tok := &token{}
	parseErr := json.Unmarshal(resBytes, tok)

	if parseErr != nil {
		return "", newJSONErr()
	}

	return tok.Access, nil
}

func (c *Client) tokensByCreds(login, pass string) (accessToken, refreshToken string, err *Error) {
	return c.tokenRequest(map[string]interface{}{
		"action":   "token",
		"login":    login,
		"password": pass,
	})
}

func (c *Client) tokenRequest(data map[string]interface{}) (accessToken, refreshToken string, err *Error) {
	_, resBytes, err := c.request(data)

	if err != nil {
		return "", "", err
	}

	tok := &token{}
	parseErr := json.Unmarshal(resBytes, tok)

	if parseErr != nil {
		return "", "", newJSONErr()
	}

	return tok.Access, tok.Refresh, nil
}
