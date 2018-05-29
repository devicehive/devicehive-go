package devicehive_go

import "github.com/devicehive/devicehive-go/internal/transport"

const (
	InvalidResponseErr  = "invalid response"
	InvalidRequestErr   = "invalid request"
	ServiceErr          = "service error"
	ConnectionFailedErr = "connection failed"
)

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

func newJSONErr() *Error {
	return &Error{name: InvalidResponseErr, reason: "data is not valid JSON string"}
}

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