// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package transport

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

type httpTimeoutErr interface {
	Timeout() bool
}

func isTimeoutErr(err error) bool {
	timeoutErr, ok := err.(httpTimeoutErr)
	return ok && timeoutErr.Timeout()
}
