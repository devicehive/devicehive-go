package transport

const (
	ConnClosedErr      = "connection closed"
	InvalidResponseErr = "invalid response"
	InvalidRequestErr  = "invalid request"
)

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