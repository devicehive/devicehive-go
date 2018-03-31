package dh

import "github.com/devicehive/devicehive-go/internal/transport"

const (
	ConnClosedErr      = transport.ConnClosedErr
	InvalidResponseErr = transport.InvalidResponseErr
	InvalidRequestErr  = transport.InvalidRequestErr
	ServiceErr         = "service error"
)

func newTransportErr(err *transport.Error) *Error {
	return &Error{name: err.Name(), reason: err.Error()}
}

type Error struct {
	name   string
	reason string
}

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
