package transport

import (
	"github.com/gorilla/websocket"
	"time"
)

const (
	DefaultTimeout = 3 * time.Second
)

type Transporter interface {
	Request(resource string, data devicehiveData, timeout time.Duration) (res []byte, err *Error)
	Subscribe(subscriptionId string) (eventChan chan []byte)
	Unsubscribe(subscriptionId string)
}

func Create(url string) (transport Transporter, err error) {
	// @TODO HTTP transport
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		return nil, err
	}

	return newWS(conn), nil
}
