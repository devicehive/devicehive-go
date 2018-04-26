package transport

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
)

func newWS(addr string) (tsp *ws, err error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)

	if err != nil {
		return nil, err
	}

	tsp = &ws{
		conn:          conn,
		requests:      make(clientsMap),
		subscriptions: make(clientsMap),
	}

	go tsp.handleServerMessages()

	return tsp, nil
}

type ws struct {
	conn          *websocket.Conn
	requests      clientsMap
	subscriptions clientsMap
}

type ids struct {
	Request      string `json:"requestId"`
	Subscription int64  `json:"subscriptionId"`
}

func (t *ws) Request(resource string, data devicehiveData, timeout time.Duration) (res []byte, err *Error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	if data == nil {
		data = devicehiveData(make(map[string]interface{}))
	}

	reqId := data.requestId()
	client := t.requests.createClient(reqId)

	data["action"] = resource

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
		return nil, &Error{name: TimeoutErr, reason: "response timeout"}
	}
}

func (t *ws) Subscribe(subscriptionId string) (eventChan chan []byte) {
	if _, ok := t.subscriptions.get(subscriptionId); ok {
		return nil
	}

	client := t.subscriptions.createSubscriber(subscriptionId)
	return client.data
}

func (t *ws) Unsubscribe(subscriptionId string) {
	client, ok := t.subscriptions.get(subscriptionId)

	if ok {
		client.close()
		t.subscriptions.delete(subscriptionId)
	}
}

func (t *ws) handleServerMessages() {
	for {
		mt, msg, err := t.conn.ReadMessage()

		connClosed := mt == websocket.CloseMessage || mt == -1
		if connClosed {
			t.terminateClients(ConnClosedErr, err)
			return
		}

		t.resolveReceiver(msg)
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

func (t *ws) resolveReceiver(msg []byte) {
	ids := &ids{}
	err := json.Unmarshal(msg, ids)

	if err != nil {
		log.Printf("request is not JSON or requestId/subscriptionId is not valid: %s", string(msg))
		return
	}

	if client, ok := t.requests.get(ids.Request); ok {
		client.data <- msg
		client.close()
		t.requests.delete(ids.Request)
	} else if client, ok := t.subscriptions.get(strconv.FormatInt(ids.Subscription, 10)); ok {
		client.data <- msg
	}
}
