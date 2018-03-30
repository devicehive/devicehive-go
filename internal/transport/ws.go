package transport

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

func newWS(conn *websocket.Conn) *ws {
	tsp := &ws{
		conn:     conn,
		requests: make(requestMap),
	}

	go tsp.handleResponses()

	return tsp
}

type ws struct {
	conn     *websocket.Conn
	requests requestMap
}

func (t *ws) Request(data devicehiveData) (res devicehiveData, err error) {
	reqId := data.requestId()
	err = t.conn.WriteJSON(data)

	if err != nil {
		return nil, err
	}

	resChan, errChan := t.requests.create(reqId)

	select {
	case res := <-resChan:
		return res, nil
	case err := <-errChan:
		return nil, err
	}
}

func (t *ws) handleResponses() {
	for {
		res := make(devicehiveData)
		mt, msg, err := t.conn.ReadMessage()

		connClosed := mt == websocket.CloseMessage || mt == -1
		if connClosed {
			t.closePendingWithErr("connection closed", err)
			return
		}

		err = json.Unmarshal(msg, &res)

		if err != nil {
			t.closePendingWithErr("invalid service response", err)
			return
		}

		t.respond(res)
	}
}

func (t *ws) closePendingWithErr(errMsg string, err error) {
	t.requests.forEach(func(resChan *response) {
		resChan.err <- fmt.Errorf("%s, error: %v", errMsg, err)
		resChan.close()
	})
}

func (t *ws) respond(res devicehiveData) {
	reqId := res.requestId()
	if resChan, ok := t.requests.get(reqId); ok {
		resChan.data <- res
		resChan.close()
		t.requests.delete(reqId)
	}
}

