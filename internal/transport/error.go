// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"
	"net/url"
	"strings"
)

const (
	ConnClosedErr      = "connection closed"
	InvalidResponseErr = "invalid response"
	InvalidRequestErr  = "invalid request"
	TimeoutErr         = "timeout"
)

func NewError(name, reason string) *Error {
	return &Error{
		name:   name,
		reason: reason,
	}
}

type Error struct {
	name   string
	reason string
}

func (e *Error) Name() string {
	return e.name
}

func (e *Error) Error() string {
	return e.reason
}

func isTimeoutErr(err error) bool {
	if err == nil {
		return false
	}

	timeoutErr, ok := err.(net.Error)
	return ok && timeoutErr.Timeout()
}

func isHostErr(err error) bool {
	if err == nil {
		return false
	}

	urlErr, ok := err.(*url.Error)
	noSuchHost := strings.HasSuffix(urlErr.Error(), "no such host")
	noData := strings.Contains(urlErr.Error(), "name is valid, but no data")

	return ok && (noSuchHost || noData)
}
