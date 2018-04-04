package transport

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
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

func (t *ws) Request(data devicehiveData, timeout time.Duration) (res []byte, err *Error) {
	reqId := data.requestId()
	wErr := t.conn.WriteJSON(data)

	if timeout == 0 {
		timeout = Timeout
	}

	if wErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: wErr.Error()}
	}

	req := t.requests.create(reqId)

	select {
	case res = <-req.response:
		return res, nil
	case err := <-req.err:
		return nil, err
	case <-time.After(timeout):
		req.close()
		t.requests.delete(reqId)
		return nil, &Error{name: TimeoutErr, reason: "request timeout"}
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
	t.requests.forEach(func(resChan *request) {
		resChan.err <- &Error{name: errMsg, reason: err.Error()}
		resChan.close()
	})
}

func (t *ws) respond(res []byte) {
	reqId := &requestId{}
	err := json.Unmarshal(res, reqId)

	if err != nil {
		log.Printf("request is not JSON or requestId is not valid: %s", string(res))
		return
	}

	if resChan, ok := t.requests.get(reqId.Value); ok {
		resChan.response <- res
		resChan.close()
		t.requests.delete(reqId.Value)
	}
}

type requestId struct {
	Value string `json:"requestId"`
}
