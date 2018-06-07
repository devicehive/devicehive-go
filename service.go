// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

import (
	"github.com/devicehive/devicehive-go/transport"
	"github.com/devicehive/devicehive-go/transportadapter"
)

// Method uses access token directly to connect
// If access token is empty it will get access token by refresh token
// It will recreate access token on expiration by given refresh token
func ConnectWithToken(url, accessToken, refreshToken string) (*Client, *Error) {
	c, err := connect(url)

	if err != nil {
		return nil, err
	}

	c.refreshToken = refreshToken

	if accessToken != "" {
		return auth(accessToken, c)
	}

	accessToken, err = c.accessTokenByRefresh(refreshToken)

	if err != nil {
		return nil, err
	}

	return auth(accessToken, c)
}

// Method obtains access token by credentials and then connects
// It will recreate access token on expiration by given credentials
func ConnectWithCreds(url, login, password string) (*Client, *Error) {
	c, err := connect(url)

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

func connect(url string) (*Client, *Error) {
	tsp, tspErr := transport.Create(url)

	if tspErr != nil {
		return nil, &Error{name: ConnectionFailedErr, reason: tspErr.Error()}
	}

	client := &Client{
		transport:                 tsp,
		transportAdapter:          transportadapter.New(tsp),
		PollingWaitTimeoutSeconds: DefaultPollingWaitTimeoutSeconds,
	}

	info, err := client.GetInfo()
	if err == nil {
		client.subscriptionTimestamp = info.ServerTimestamp.Time
	}

	return client, nil
}

func auth(accTok string, c *Client) (*Client, *Error) {
	auth, err := c.authenticate(accTok)

	if err != nil {
		return nil, err
	}

	if auth {
		return c, nil
	}

	return nil, nil
}
