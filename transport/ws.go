package transport

import (
	"github.com/devicehive/devicehive-go/transport/apirequests"
	"github.com/devicehive/devicehive-go/utils"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
	"sync"
)

func newWS(addr string) (tsp *ws, err error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)

	if err != nil {
		return nil, err
	}

	tsp = &ws{
		conn:          conn,
		connMu:		   sync.Mutex{},
		requests:      apirequests.NewClientsMap(),
		subscriptions: apirequests.NewWSSubscriptionsMap(apirequests.NewClientsMap()),
	}

	go tsp.handleServerMessages()

	return tsp, nil
}

type ws struct {
	conn          *websocket.Conn
	connMu		  sync.Mutex
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
	req := t.requests.CreateRequest(reqId)

	data := params.mapData()
	data["action"] = resource
	data["requestId"] = reqId

	t.connMu.Lock()
	wErr := t.conn.WriteJSON(data)
	t.connMu.Unlock()
	if wErr != nil {
		return nil, NewError(InvalidRequestErr, wErr.Error())
	}

	select {
	case res = <-req.Data:
		return res, nil
	case err := <-req.Err:
		return nil, NewError(ConnClosedErr, err.Error())
	case <-time.After(timeout):
		req.Close()
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

	subscription := t.subscriptions.CreateSubscription(subscriptionId)
	return subscription.Data
}

func (t *ws) Unsubscribe(subscriptionId string) {
	subscription, ok := t.subscriptions.Get(subscriptionId)

	if ok {
		subscription.Close()
		t.subscriptions.Delete(subscriptionId)
	}
}

func (t *ws) handleServerMessages() {
	for {
		mt, msg, err := t.conn.ReadMessage()

		connClosed := mt == websocket.CloseMessage || mt == -1
		if connClosed {
			t.terminateRequests(err)
			return
		}

		t.resolveReceiver(msg)
	}
}

func (t *ws) terminateRequests(err error) {
	t.requests.ForEach(func(req *apirequests.PendingRequest) {
		req.Err <- err
		req.Close()
	})

	t.subscriptions.ForEach(func(req *apirequests.PendingRequest) {
		req.Close()
	})
}

func (t *ws) resolveReceiver(msg []byte) {
	ids, err := utils.ParseIDs(msg)

	if err != nil {
		log.Printf("request is not JSON or requestId/subscriptionId is not valid: %s", string(msg))
		return
	}

	if req, ok := t.requests.Get(ids.Request); ok {
		req.Data <- msg
		req.Close()
		t.requests.Delete(ids.Request)
	} else if ids.Subscription != 0 {
		subsId := strconv.FormatInt(ids.Subscription, 10)
		if subs, ok := t.subscriptions.Get(subsId); ok {
			subs.Data <- msg
		} else {
			t.subscriptions.BufferPut(msg)
		}
	}
}
