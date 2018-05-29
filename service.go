package devicehive_go

import (
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/transportadapter"
	"time"
)

const (
	NotificationType = "notification"
	CommandType      = "command"
	Timeout          = 5 * time.Second
)

var client = &Client{
	PollingWaitTimeoutSeconds: 30,
}

func ConnectWithToken(url, accessToken, refreshToken string) (c *Client, err *Error) {
	client, err = connect(url)

	if err != nil {
		return nil, err
	}

	client.refreshToken = refreshToken

	return auth(accessToken, client)
}

func ConnectWithCreds(url, login, password string) (c *Client, err *Error) {
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

func connect(url string) (c *Client, err *Error) {
	tsp, tspErr := transport.Create(url)

	if tspErr != nil {
		return nil, &Error{name: ConnectionFailedErr, reason: tspErr.Error()}
	}

	client.transport = tsp
	client.transportAdapter = transportadapter.New(client.transport)

	return client, nil
}

func auth(accTok string, c *Client) (client *Client, err *Error) {
	auth, err := c.authenticate(accTok)

	if err != nil {
		return nil, err
	}

	if auth {
		return c, nil
	}

	return nil, nil
}