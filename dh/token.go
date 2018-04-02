package dh

import (
	"time"
	"encoding/json"
)

type token struct {
	Access string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}

func (c *Client) TokenByCreds(login, pass string) (accessToken, refreshToken string, err *Error) {
	return c.tokenRequest(map[string]interface{}{
		"action":   "token",
		"login":    login,
		"password": pass,
	})
}

func (c *Client) TokenByPayload(userId int, actions, networkIds, deviceTypeIds []string, expiration *time.Time) (accessToken, refreshToken string, err *Error) {
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
	if expiration != nil {
		payload["expiration"] = expiration.UTC().Format(timestampLayout)
	}

	data := map[string]interface{}{
		"action":  "token/create",
		"payload": payload,
	}

	return c.tokenRequest(data)
}

func (c *Client) tokenRequest(data map[string]interface{}) (accessToken, refreshToken string, err *Error) {
	resBytes, tspErr := c.tsp.Request(data)

	if _, err = c.handleResponse(resBytes, tspErr); err != nil {
		return "", "", err
	}

	tok := &token{}
	parseErr := json.Unmarshal(resBytes, tok)

	if parseErr != nil {
		return "", "", newJSONErr()
	}

	return tok.Access, tok.Refresh, nil
}

func (c *Client) TokenRefresh(refreshToken string) (accessToken string, err *Error) {
	resBytes, tspErr := c.tsp.Request(map[string]interface{}{
		"action":       "token/refresh",
		"refreshToken": refreshToken,
	})

	if _, err = c.handleResponse(resBytes, tspErr); err != nil {
		return "", err
	}

	tok := &token{}
	parseErr := json.Unmarshal(resBytes, tok)

	if parseErr != nil {
		return "", newJSONErr()
	}

	return tok.Access, nil
}
