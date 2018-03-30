package dh

import (
	"github.com/devicehive/devicehive-go/internal/transport"
)

func Connect(url string) (*Client, error) {
	tsp, err := transport.Create(url)

	if err != nil {
		return nil, err
	}

	return &Client{tsp: tsp}, nil
}
