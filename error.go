// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

import "github.com/devicehive/devicehive-go/transport"

func newError(err error) *Error {
	if err == nil {
		return nil
	}

	switch err.(type) {
	case *transport.Error:
		return newTransportErr(err.(*transport.Error))
	default:
		return &Error{ServiceErr, err.Error()}
	}
}

func newJSONErr(err error) *Error {
	return &Error{name: InvalidResponseErr, reason: "data is not valid JSON string: " + err.Error()}
}

func newTransportErr(err *transport.Error) *Error {
	return &Error{name: err.Name(), reason: err.Error()}
}

type Error struct {
	name   string
	reason string
}

// Method serves as name getter to classify the type of error
func (e *Error) Name() string {
	if e == nil {
		return ""
	}

	return e.name
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	return e.reason
}
