package dh

import (
	"github.com/devicehive/devicehive-go/internal/transport"
	"time"
)

const (
	NotificationType = "notification"
	CommandType      = "command"
	timestampLayout  = "2006-01-02T15:04:05.000"
	Timeout          = 1 * time.Second
)

func Connect(url string) (client *Client, err *Error) {
	tsp, tspErr := transport.Create(url)

	if tspErr != nil {
		return nil, &Error{name: ConnectionFailedErr, reason: tspErr.Error()}
	}

	return &Client{tsp: tsp}, nil
}
