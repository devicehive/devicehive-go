package dh

import (
	"github.com/gorilla/websocket"
)

func Connect(url string) (*dhClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		return nil, err
	}

	return &dhClient{conn: conn}, nil
}