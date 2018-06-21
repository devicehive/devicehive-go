// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package transport

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/devicehive/devicehive-go/internal/transport/apirequests"
	"github.com/devicehive/devicehive-go/internal/utils"
	"github.com/gorilla/websocket"
)

func newWS(addr string, params *Params) (tsp *WS, err error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	tsp = &WS{
		conn:          conn,
		address:       addr,
		requests:      apirequests.NewClientsMap(),
		subscriptions: apirequests.NewWSSubscriptionsMap(apirequests.NewClientsMap()),
	}
	tsp.setParams(params)

	go tsp.handleServerMessages()

	return tsp, nil
}

type WS struct {
	conn                 *websocket.Conn
	address              string
	connMu               sync.Mutex
	requests             *apirequests.PendingRequestsMap
	subscriptions        *apirequests.WSSubscriptionsMap
	reconnectionTries    int
	reconnectionInterval time.Duration
}

func (t *WS) IsHTTP() bool {
	return false
}

func (t *WS) IsWS() bool {
	return true
}

func (t *WS) Request(resource string, params *RequestParams, timeout time.Duration) ([]byte, *Error) {
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
	case res := <-req.Data:
		return res, nil
	case err := <-req.Err:
		return nil, NewError(ConnClosedErr, err.Error())
	case <-time.After(timeout):
		req.Close()
		t.requests.Delete(reqId)
		return nil, NewError(TimeoutErr, "response timeout")
	}
}

func (t *WS) Subscribe(resource string, params *RequestParams) (subscription *Subscription, subscriptionId string, err *Error) {
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

func (t *WS) subscribe(subscriptionId string) *Subscription {
	if _, ok := t.subscriptions.Get(subscriptionId); ok {
		return nil
	}

	subs := t.subscriptions.CreateSubscription(subscriptionId)

	subscription := &Subscription{
		DataChan: subs.Data,
		ErrChan:  subs.Err,
	}

	return subscription
}

func (t *WS) Unsubscribe(subscriptionId string) {
	subscription, ok := t.subscriptions.Get(subscriptionId)

	if ok {
		subscription.Close()
		t.subscriptions.Delete(subscriptionId)
	}
}

func (t *WS) handleServerMessages() {
	for {
		mt, msg, err := t.conn.ReadMessage()
		if mt == websocket.CloseMessage {
			t.terminateRequests(err)
			return
		}

		serverDown := mt == -1
		if serverDown {
			reconnectDisabled := t.reconnectionTries == 0 || t.reconnectionInterval == 0
			if reconnectDisabled {
				t.terminateRequests(err)
				return
			}

			reconnErr := t.reconnect()
			if reconnErr != nil {
				t.terminateRequests(reconnErr)
				return
			}

			continue
		}

		t.resolveReceiver(msg)
	}
}

func (t *WS) reconnect() error {
	var reconnErr error
	for i := 0; i < t.reconnectionTries; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(t.address, nil)
		if err != nil {
			reconnErr = err
			time.Sleep(t.reconnectionInterval)
			continue
		}

		t.conn = conn
		reconnErr = nil
		break
	}

	return reconnErr
}

func (t *WS) terminateRequests(err error) {
	t.requests.ForEach(func(req *apirequests.PendingRequest) {
		req.Err <- err
		req.Close()
	})

	t.subscriptions.ForEach(func(req *apirequests.PendingRequest) {
		req.Close()
	})
}

func (t *WS) resolveReceiver(msg []byte) {
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

func (t *WS) setParams(p *Params) {
	if p == nil {
		return
	}

	if p.ReconnectionTries != 0 {
		t.reconnectionTries = p.ReconnectionTries
	}

	if p.ReconnectionInterval != 0 {
		t.reconnectionInterval = p.ReconnectionInterval
	}
}
