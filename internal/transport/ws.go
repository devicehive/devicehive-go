package transport

import (
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
	"github.com/devicehive/devicehive-go/internal/transport/apirequests"
	"github.com/devicehive/devicehive-go/internal/utils"
)

func newWS(addr string) (tsp *ws, err error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)

	if err != nil {
		return nil, err
	}

	tsp = &ws{
		conn:          conn,
		requests:      apirequests.NewClientsMap(),
		subscriptions: apirequests.NewWSSubscriptionsMap(apirequests.NewClientsMap()),
	}

	go tsp.handleServerMessages()
	go tsp.subscriptions.CleanupBufferByTimeout(1 * time.Second)

	return tsp, nil
}

type ws struct {
	conn          *websocket.Conn
	requests      *apirequests.PendingRequestsMap
	subscriptions *apirequests.WSSubscriptionsMap
}

func (t *ws) IsHTTP() bool {
	return false
}

func (t *ws) IsWS() bool {
	return true
}

func (t *ws) Request(resource string, params *RequestParams, timeout time.Duration) (res []byte, err *Error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	if params == nil {
		params = &RequestParams{}
	}

	reqId := params.requestId()
	client := t.requests.CreateClient(reqId)

	data := params.mapData()
	data["action"] = resource
	data["requestId"] = reqId

	wErr := t.conn.WriteJSON(data)
	if wErr != nil {
		return nil, NewError(InvalidRequestErr, wErr.Error())
	}

	select {
	case res = <-client.Data:
		return res, nil
	case err := <-client.Err:
		return nil, NewError(ConnClosedErr, err.Error())
	case <-time.After(timeout):
		client.Close()
		t.requests.Delete(reqId)
		return nil, NewError(TimeoutErr, "response timeout")
	}
}

func (t *ws) Subscribe(resource string, params *RequestParams) (eventChan chan []byte, subscriptionId string, err *Error) {
	res, err := t.Request(resource, params, 0)
	if err != nil {
		return nil, "", err
	}

	ids, parseErr := utils.ParseIDs(res)
	if parseErr != nil {
		return nil, "", NewError(InvalidResponseErr, parseErr.Error())
	}
	subscriptionId = strconv.FormatInt(ids.Subscription, 10)

	return t.subscribe(subscriptionId), subscriptionId, nil
}

func (t *ws) subscribe(subscriptionId string) (eventChan chan []byte) {
	if _, ok := t.subscriptions.Get(subscriptionId); ok {
		return nil
	}

	client := t.subscriptions.CreateSubscriber(subscriptionId)
	return client.Data
}

func (t *ws) Unsubscribe(subscriptionId string) {
	client, ok := t.subscriptions.Get(subscriptionId)

	if ok {
		client.Close()
		t.subscriptions.Delete(subscriptionId)
	}
}

func (t *ws) handleServerMessages() {
	for {
		mt, msg, err := t.conn.ReadMessage()

		connClosed := mt == websocket.CloseMessage || mt == -1
		if connClosed {
			t.terminateClients(err)
			return
		}

		t.resolveReceiver(msg)
	}
}

func (t *ws) terminateClients(err error) {
	t.requests.ForEach(func(c *apirequests.PendingRequest) {
		c.Err <- err
		c.Close()
	})

	t.subscriptions.ForEach(func(c *apirequests.PendingRequest) {
		c.Close()
	})
}

func (t *ws) resolveReceiver(msg []byte) {
	ids, err := utils.ParseIDs(msg)

	if err != nil {
		log.Printf("request is not JSON or requestId/subscriptionId is not valid: %s", string(msg))
		return
	}

	if client, ok := t.requests.Get(ids.Request); ok {
		client.Data <- msg
		client.Close()
		t.requests.Delete(ids.Request)
	} else {
		if client, ok := t.subscriptions.Get(strconv.FormatInt(ids.Subscription, 10)); ok {
			client.Data <- msg
		} else {
			t.subscriptions.BufferPut(msg)
		}
	}
}
