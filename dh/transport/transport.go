package transport

import "github.com/gorilla/websocket"

type Transporter interface {
	Request(data request) (res response, err error)
}

type response map[string]interface{}
type request map[string]interface{}

func Create(url string) (transport Transporter, err error) {
	// @TODO HTTP transport
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		return nil, err
	}

	return newWS(conn), nil
}