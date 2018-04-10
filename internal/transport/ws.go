package transport

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
	"strconv"
)

func newWS(conn *websocket.Conn) *ws {
	tsp := &ws{
		conn:     conn,
		requests: make(requestMap),
		subscriptions: make(requestMap),
	}

	go tsp.handleResponses()

	return tsp
}

type ws struct {
	conn     *websocket.Conn
	requests requestMap
	subscriptions requestMap
}

func (t *ws) Subscribe(subscriptionId int64) (eventChan chan []byte) {
	subscription := t.subscriptions.create(strconv.FormatInt(subscriptionId, 10))
	return subscription.response
}

func (t *ws) Request(data devicehiveData, timeout time.Duration) (res []byte, err *Error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	reqId := data.requestId()
	req := t.requests.create(reqId)

	wErr := t.conn.WriteJSON(data)
	if wErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: wErr.Error()}
	}

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
	tspErr := &Error{name: errMsg, reason: err.Error()}
	closeChans := func(resChan *request) {
		resChan.err <- tspErr
		resChan.close()
	}

	t.requests.forEach(closeChans)
	t.subscriptions.forEach(closeChans)
}

func (t *ws) respond(res []byte) {
	ids := &ids{}
	err := json.Unmarshal(res, ids)

	if err != nil {
		log.Printf("request is not JSON or requestId/subscriptionId is not valid: %s", string(res))
		return
	}

	if resChan, ok := t.requests.get(ids.Request); ok {
		resChan.response <- res
		resChan.close()
		t.requests.delete(ids.Request)
	} else if eventChan, ok := t.subscriptions.get(strconv.FormatInt(ids.Subscription, 10)); ok {
		eventChan.response <- res
	}
}

type ids struct {
	Request string `json:"requestId"`
	Subscription int64 `json:"subscriptionId"`
}
