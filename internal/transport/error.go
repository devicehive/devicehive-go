package transport

const (
	ConnClosedErr      = "connection closed"
	InvalidResponseErr = "invalid request"
	InvalidRequestErr  = "invalid request"
	TimeoutErr         = "timeout"
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
