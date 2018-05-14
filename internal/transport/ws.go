package transport

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
	"fmt"
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

func (t *ws) IsHTTP() bool {
	return false
}

func (t *ws) IsWS() bool {
	return true
}

type ids struct {
	Request      string `json:"requestId"`
	Subscription int64  `json:"subscriptionId"`
}

func (t *ws) Request(resource string, params *RequestParams, timeout time.Duration) (res []byte, err *Error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	if params == nil {
		params = &RequestParams{}
	}

	reqId := params.requestId()
	client := t.requests.createClient(reqId)

	data := params.mapData()
	data["action"] = resource
	data["requestId"] = reqId

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

func (t *ws) Subscribe(resource string, params *RequestParams) (eventChan chan []byte, subscriptionId string, err *Error) {
	res, err := t.Request(resource, params, 0)
	if err != nil {
		return nil, "", err
	}

	id := &ids{}

	parseErr := json.Unmarshal(res, id)
	if parseErr != nil {
		return nil, "", &Error{name: InvalidResponseErr, reason: parseErr.Error()}
	}
	subscriptionId = strconv.FormatInt(id.Subscription, 10)

	fmt.Println("Subscribed ", subscriptionId)

	return t.subscribe(subscriptionId), subscriptionId, nil
}

func (t *ws) subscribe(subscriptionId string) (eventChan chan []byte) {
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

	fmt.Println(ids.Request, ids.Subscription)
	fmt.Println(string(msg))

	if err != nil {
		log.Printf("request is not JSON or requestId/subscriptionId is not valid: %s", string(msg))
		return
	}

	if client, ok := t.requests.get(ids.Request); ok {
		client.data <- msg
		client.close()
		t.requests.delete(ids.Request)
	} else if client, ok := t.subscriptions.get(strconv.FormatInt(ids.Subscription, 10)); ok {
		fmt.Println("Sending to data chan...")
		client.data <- msg
	}
}
