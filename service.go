package devicehive_go

import (
	"github.com/devicehive/devicehive-go/transport"
	"github.com/devicehive/devicehive-go/transportadapter"
)

// Method uses access token directly to connect
// It will recreate access token on expiration by given refresh token
func ConnectWithToken(url, accessToken, refreshToken string) (c *Client, err *Error) {
	c, err = connect(url)

	if err != nil {
		return nil, err
	}

	c.refreshToken = refreshToken

	return auth(accessToken, c)
}

// Method obtains access token by credentials and then connects
// It will recreate access token on expiration by given credentials
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
