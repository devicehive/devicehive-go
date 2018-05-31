package devicehive_go

import (
	"github.com/devicehive/devicehive-go/transport"
	"github.com/devicehive/devicehive-go/transportadapter"
)

var client = &Client{
	PollingWaitTimeoutSeconds: 30,
}

func ConnectWithToken(url, accessToken, refreshToken string) (c *Client, err *Error) {
	c, err = connect(url)

	if err != nil {
		return nil, err
	}

	c.refreshToken = refreshToken

	return auth(accessToken, c)
}

func ConnectWithCreds(url, login, password string) (c *Client, err *Error) {
	c, err = connect(url)

	if err != nil {
		return nil, err
	}

	accTok, _, err := c.tokensByCreds(login, password)

	if err != nil {
		return nil, err
	}

	c.login = login
	c.password = password

	return auth(accTok, c)
}

func connect(url string) (c *Client, err *Error) {
	tsp, tspErr := transport.Create(url)

	if tspErr != nil {
		return nil, &Error{name: ConnectionFailedErr, reason: tspErr.Error()}
	}

	client := &Client{
		transport:                 tsp,
		transportAdapter:          transportadapter.New(tsp),
		PollingWaitTimeoutSeconds: DefaultPollingWaitTimeoutSeconds,
	}

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
