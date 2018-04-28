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
		"payload": payload,
	}

	return c.tokenRequest("tokenCreate", data)
}

func (c *Client) RefreshToken() (accessToken string, err *Error) {
	if c.refreshToken == "" {
		accessToken, _, err = c.tokensByCreds(c.login, c.password)
		return accessToken, err
	}

	return c.accessTokenByRefresh(c.refreshToken)
}

func (c *Client) accessTokenByRefresh(refreshToken string) (accessToken string, err *Error) {
	_, resBytes, err := c.request("token/refresh", map[string]interface{}{
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
	return c.tokenRequest("token", map[string]interface{}{
		"login":    login,
		"password": pass,
	})
}

func (c *Client) tokenRequest(resourceName string, data map[string]interface{}) (accessToken, refreshToken string, err *Error) {
	_, resBytes, err := c.request(resourceName, data)

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
