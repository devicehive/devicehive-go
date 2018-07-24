// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package transportadapter

import (
	"time"

	"github.com/devicehive/devicehive-go/internal/transport"
)

func New(tsp transport.Transporter) TransportAdapter {
	if tsp, ok := tsp.(*transport.WS); ok {
		ws := newWSAdapter(tsp)
		return ws
	}

	http := &HTTPAdapter{
		transport: tsp.(*transport.HTTP),
	}
	return http
}

type TransportAdapter interface {
	Request(resourceName string, data map[string]interface{}, timeout time.Duration) (res []byte, err error)
	Subscribe(resourceName string, pollingWaitTimeoutSeconds int, params map[string]interface{}) (subscription *transport.Subscription, subscriptionId string, err *transport.Error)
	Unsubscribe(resourceName, subscriptionId string, timeout time.Duration) error
	Authenticate(token string, timeout time.Duration) (result bool, err error)
	SetCreds(login, password string)
	SetRefreshToken(refTok string)
	RefreshToken() (accessToken string, err error)
	TokensByCreds(login, pass string) (accessToken, refreshToken string, err error)
	AccessTokenByRefresh(refreshToken string) (accessToken string, err error)
}
