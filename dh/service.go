package dh

import (
	"github.com/devicehive/devicehive-go/internal/transport"
	"time"
)

const (
	NotificationType = "notification"
	Timeout          = 5 * time.Second
)

func ConnectWithToken(url, accessToken, refreshToken string) (client *Client, err *Error) {
	client, err = connect(url)

	if err != nil {
		return nil, err
	}

	client.refreshToken = refreshToken

	return auth(accessToken, client)
}

func ConnectWithCreds(url, login, password string) (client *Client, err *Error) {
	client, err = connect(url)

	if err != nil {
		return nil, err
	}

	accTok, _, err := client.tokensByCreds(login, password)

	if err != nil {
		return nil, err
	}

	client.login = login
	client.password = password

	return auth(accTok, client)
}

func connect(url string) (client *Client, err *Error) {
	tsp, tspErr := transport.Create(url)

	if tspErr != nil {
		return nil, &Error{name: ConnectionFailedErr, reason: tspErr.Error()}
	}

	return &Client{tsp: tsp}, nil
}

func auth(accTok string, c *Client) (client *Client, err *Error) {
	auth, err := c.Authenticate(accTok)

	if err != nil {
		return nil, err
	}

	if auth {
		return c, nil
	}

	return nil, nil
}
