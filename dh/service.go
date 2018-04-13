package dh

import (
	"github.com/devicehive/devicehive-go/internal/transport"
	"time"
)

const (
	NotificationType = "notification"
	Timeout          = 5 * time.Second
)

func Connect(url string) (client *Client, err *Error) {
	tsp, tspErr := transport.Create(url)

	if tspErr != nil {
		return nil, &Error{name: ConnectionFailedErr, reason: tspErr.Error()}
	}

	return &Client{tsp: tsp}, nil
}
