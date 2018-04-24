package transport

const (
	ConnClosedErr      = "connection closed"
	InvalidResponseErr = "invalid response"
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

type httpTimeoutErr interface {
	Timeout() bool
}

func isTimeoutErr(err error) bool {
	timeoutErr, ok := err.(httpTimeoutErr)
	return ok && timeoutErr.Timeout()
}
