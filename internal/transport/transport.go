// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package transport

import (
	"github.com/devicehive/devicehive-go/internal/requestparams"
	"net/url"
	"time"
)

const (
	DefaultTimeout = 5 * time.Second
)

type Transporter interface {
	Request(resource string, params *requestparams.RequestParams, timeout time.Duration) (res []byte, err *Error)
	Subscribe(resource string, params *requestparams.RequestParams) (subscription *Subscription, subscriptionId string, err *Error)
	Unsubscribe(subscriptionId string)
}

func Create(addr string, p *Params) (Transporter, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "http" || u.Scheme == "https" {
		return newHTTP(addr, p)
	}

	return newWS(addr, p)
}
