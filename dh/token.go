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
	data := map[string]interface{}{
		"userId": userId,
	}

	if actions != nil {
		data["actions"] = actions
	}
	if networkIds != nil {
		data["networkIds"] = networkIds
	}
	if deviceTypeIds != nil {
		data["deviceTypeIds"] = deviceTypeIds
	}
	if expiration.Unix() > 0 {
		data["expiration"] = expiration.UTC().Format(timestampLayout)
	}

	return c.tokenRequest("tokenCreate", map[string]interface{} {
		"payload": data,
	})
}

func (c *Client) RefreshToken() (accessToken string, err *Error) {
	if c.refreshToken == "" {
		accessToken, _, err = c.tokensByCreds(c.login, c.password)
		return accessToken, err
	}

	return c.accessTokenByRefresh(c.refreshToken)
}

func (c *Client) accessTokenByRefresh(refreshToken string) (accessToken string, err *Error) {
	rawRes, err := c.request("tokenRefresh", map[string]interface{}{
		"refreshToken": c.refreshToken,
	})

	if err != nil {
		return "", err
	}

	tok := &token{}
	parseErr := json.Unmarshal(rawRes, tok)

	if parseErr != nil {
		return "", newJSONErr()
	}

	return tok.Access, nil
}

func (c *Client) tokensByCreds(login, pass string) (accessToken, refreshToken string, err *Error) {
	return c.tokenRequest("tokenByCreds", map[string]interface{}{
		"login":    login,
		"password": pass,
	})
}

func (c *Client) tokenRequest(resourceName string, data map[string]interface{}) (accessToken, refreshToken string, err *Error) {
	rawRes, err := c.request(resourceName, data)

	if err != nil {
		return "", "", err
	}

	tok := &token{}
	parseErr := json.Unmarshal(rawRes, tok)

	if parseErr != nil {
		return "", "", newJSONErr()
	}

	return tok.Access, tok.Refresh, nil
}
