// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

import (
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/internal/transportadapter"
)

// Method uses access token directly to connect
// If access token is empty it will get access token by refresh token
// It will recreate access token on expiration by given refresh token
func ConnectWithToken(url, accessToken, refreshToken string, p *ConnectionParams) (*Client, *Error) {
	c, err := connect(url, p)
	if err != nil {
		return nil, err
	}

	c.setRefreshToken(refreshToken)

	if accessToken != "" {
		return auth(accessToken, c)
	}

	accessToken, err = c.RefreshToken()

	if err != nil {
		return nil, err
	}

	return auth(accessToken, c)
}

// Method obtains access token by credentials and then connects
// It will recreate access token on expiration by given credentials
func ConnectWithCreds(url, login, password string, p *ConnectionParams) (*Client, *Error) {
	c, err := connect(url, p)
	if err != nil {
		return nil, err
	}

	c.setCreds(login, password)

	accTok, err := c.RefreshToken()

	if err != nil {
		return nil, err
	}

	return auth(accTok, c)
}

func connect(url string, p *ConnectionParams) (*Client, *Error) {
	timeout := p.Timeout()
	tspParams := createTransportParams(p)

	tsp, tspErr := transport.Create(url, tspParams)

	if tspErr != nil {
		return nil, &Error{name: ConnectionFailedErr, reason: tspErr.Error()}
	}

	client := &Client{
		transportAdapter:          transportadapter.New(tsp),
		PollingWaitTimeoutSeconds: DefaultPollingWaitTimeoutSeconds,
		defaultRequestTimeout:     timeout,
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

func createTransportParams(p *ConnectionParams) *transport.Params {
	var tspParams *transport.Params
	if p != nil {
		tspParams = &transport.Params{
			ReconnectionTries:    p.ReconnectionTries,
			ReconnectionInterval: p.ReconnectionInterval,
			DefaultTimeout:       p.Timeout(),
		}
	}

	return tspParams
}
