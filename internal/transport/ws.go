package transport

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
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

func (t *ws) Request(data devicehiveData) (res []byte, err *Error) {
	reqId := data.requestId()
	wErr := t.conn.WriteJSON(data)

	if wErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: wErr.Error()}
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
		mt, msg, err := t.conn.ReadMessage()

		connClosed := mt == websocket.CloseMessage || mt == -1
		if connClosed {
			t.closePendingWithErr(ConnClosedErr, err)
			return
		}

		t.respond(msg)
	}
}

func (t *ws) closePendingWithErr(errMsg string, err error) {
	t.requests.forEach(func(resChan *response) {
		resChan.err <- &Error{name: errMsg, reason: err.Error()}
		resChan.close()
	})
}

func (t *ws) respond(res []byte) {
	reqId := &requestId{}
	err := json.Unmarshal(res, reqId)

	if err != nil {
		log.Printf("response is not JSON or requestId is not valid: %s", string(res))
		return
	}

	if resChan, ok := t.requests.get(reqId.Value); ok {
		resChan.data <- res
		resChan.close()
		t.requests.delete(reqId.Value)
	}
}

type requestId struct {
	Value string `json:"requestId"`
}
