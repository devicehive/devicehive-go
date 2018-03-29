package transport

import (
	"github.com/gorilla/websocket"
)

func newWS(conn *websocket.Conn) *ws {
	tsp := &ws{
		conn: conn,
		requests: make(requestMap),
	}

	go tsp.handleResponses()

	return tsp
}

type ws struct {
	conn *websocket.Conn
	requests requestMap
}

func (t *ws) Request(data devicehiveData) (res devicehiveData, err error) {
	reqId := data.requestId()
	err = t.conn.WriteJSON(data)

	if err != nil {
		return nil, err
	}

	t.requests.add(reqId, make(chan devicehiveData))

	res = <- t.requests[reqId]

	return res, nil
}

func (t *ws) handleResponses() {
	for {
		res := make(devicehiveData)
		err := t.conn.ReadJSON(&res)

		// @TODO ReadJSON error handling
		if err != nil {
			panic(err)
		}

		reqId := res["requestId"].(string)
		if resChan, ok := t.requests.get(reqId); ok {
			resChan <- res
			close(resChan)
			t.requests.delete(reqId)
		}
	}
}