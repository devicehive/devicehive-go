package dh

import (
	"github.com/devicehive/devicehive-go/dh/transport"
	"math/rand"
	"time"
	"fmt"
	"strings"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Client struct {
	tsp transport.Transporter
}

func (c *Client) Authenticate(token string) (result bool, err error) {
	res, err := c.tsp.Request(map[string]interface{}{
		"action": "authenticate",
		"token":  token,
	})

	if err != nil {
		return false, err
	}

	return res["status"].(string) == "success", nil
}

func (c *Client) TokenByCreds(login, pass string) (accessToken, refreshToken string, err error) {
	return c.tokenRequest(map[string]interface{}{
		"action":   "token",
		"login":    login,
		"password": pass,
	})
}

func (c *Client) TokenByPayload(userId int, actions, networkIds, deviceTypeIds []string, expiration time.Time) (accessToken, refreshToken string, err error) {
	return c.tokenRequest(map[string]interface{}{
		"action": "token/create",
		"payload": map[string]interface{}{
			"userId": userId,
			"actions": actions,
			"networkIds": networkIds,
			"deviceTypeIds": deviceTypeIds,
			"expiration": expiration.String(),
		},
	})
}

func (c *Client) tokenRequest(data map[string]interface{}) (accessToken, refreshToken string, err error) {
	res, err := c.tsp.Request(data)

	if err != nil {
		return "", "", err
	}

	if res["status"] == "error" {
		errMsg := strings.ToLower(res["error"].(string))
		errCode := int(res["code"].(float64))
		return "", "", fmt.Errorf("%d %s", errCode, errMsg)
	}

	return res["accessToken"].(string), res["refreshToken"].(string), nil
}

func (c *Client) TokenRefresh(refreshToken string) (accessToken string, err error) {
	res, err := c.tsp.Request(map[string]interface{}{
		"action":       "token/refresh",
		"refreshToken": refreshToken,
	})

	if err != nil {
		return "", err
	}

	return res["accessToken"].(string), nil
}
