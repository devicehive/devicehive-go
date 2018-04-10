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
		requests: make(clientsMap),
		subscriptions: make(clientsMap),
	}

	go tsp.handleResponses()

	return tsp
}

type ws struct {
	conn          *websocket.Conn
	requests      clientsMap
	subscriptions clientsMap
}

func (t *ws) Subscribe(subscriptionId string) (eventChan chan []byte) {
	client := t.subscriptions.createSubscriber(subscriptionId)
	return client.data
}

func (t *ws) Request(data devicehiveData, timeout time.Duration) (res []byte, err *Error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	reqId := data.requestId()
	client := t.requests.createClient(reqId)

	wErr := t.conn.WriteJSON(data)
	if wErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: wErr.Error()}
	}

	select {
	case res = <-client.data:
		return res, nil
	case err := <-client.err:
		return nil, err
	case <-time.After(timeout):
		client.close()
		t.requests.delete(reqId)
		return nil, &Error{name: TimeoutErr, reason: "request timeout"}
	}
}

func (t *ws) handleResponses() {
	for {
		mt, msg, err := t.conn.ReadMessage()

		connClosed := mt == websocket.CloseMessage || mt == -1
		if connClosed {
			t.terminateClients(ConnClosedErr, err)
			return
		}

		t.respond(msg)
	}
}

func (t *ws) terminateClients(errMsg string, err error) {
	tspErr := &Error{name: errMsg, reason: err.Error()}
	t.requests.forEach(func(c *client) {
		c.err <- tspErr
		c.close()
	})

	t.subscriptions.forEach(func(c *client) {
		c.close()
	})
}

func (t *ws) respond(res []byte) {
	ids := &ids{}
	err := json.Unmarshal(res, ids)

	if err != nil {
		log.Printf("request is not JSON or requestId/subscriptionId is not valid: %s", string(res))
		return
	}

	if client, ok := t.requests.get(ids.Request); ok {
		client.data <- res
		client.close()
		t.requests.delete(ids.Request)
	} else if client, ok := t.subscriptions.get(strconv.FormatInt(ids.Subscription, 10)); ok {
		client.data <- res
	}
}

type ids struct {
	Request string `json:"requestId"`
	Subscription int64 `json:"subscriptionId"`
}
