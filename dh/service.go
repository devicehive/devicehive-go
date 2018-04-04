package dh

import (
	"github.com/devicehive/devicehive-go/internal/transport"
	"time"
)

const (
	Notification    = "notification"
	Command         = "command"
	timestampLayout = "2006-01-02T15:04:05.000"
	Timeout         = 1 * time.Second
)

func Connect(url string) (*Client, error) {
	tsp, err := transport.Create(url)

	if err != nil {
		return nil, err
	}

	return &Client{tsp: tsp}, nil
}
