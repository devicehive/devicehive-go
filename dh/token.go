package dh

import "time"

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
	res, tspErr := c.tsp.Request(data)

	if err = c.handleResponse(res, tspErr); err != nil {
		return "", "", err
	}

	return res["accessToken"].(string), res["refreshToken"].(string), nil
}

func (c *Client) TokenRefresh(refreshToken string) (accessToken string, err *Error) {
	res, tspErr := c.tsp.Request(map[string]interface{}{
		"action":       "token/refresh",
		"refreshToken": refreshToken,
	})

	if err = c.handleResponse(res, tspErr); err != nil {
		return "", err
	}

	return res["accessToken"].(string), nil
}
