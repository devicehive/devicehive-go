package dh

import (
	"math/rand"
	"time"
	"github.com/devicehive/devicehive-go/dh/transport"
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
		"token": token,
	})

	if err != nil {
		return false, err
	}

	return res["status"].(string) == "success", nil
}

func (c *Client) TokenByCreds(login, pass string) (accessToken, refreshToken string, err error) {
	res, err := c.tsp.Request(map[string]interface{}{
		"action": "token",
		"login": login,
		"password": pass,
	})

	if err != nil {
		return "", "", err
	}

	return res["accessToken"].(string), res["refreshToken"].(string), nil
}

func (c *Client) TokenRefresh(refreshToken string) (accessToken string, err error) {
	res, err := c.tsp.Request(map[string]interface{}{
		"action": "token/refresh",
		"refreshToken": refreshToken,
	})

	if err != nil {
		return "", err
	}

	return res["accessToken"].(string), nil
}