package transportadapter

import (
	"fmt"
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/internal/transportadapter"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"testing"
	"time"
	"sync"
)

func TestWSReconnection(t *testing.T) {
	var srv *http.Server
	var conn *websocket.Conn
	var mu sync.Mutex
	port := rand.Int31n(45535) + 20000
	addr := fmt.Sprintf("localhost:%d", port)

	go func() {
		srv = &http.Server{
			Addr: addr,
		}

		u := &websocket.Upgrader{}
		srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := u.Upgrade(w, r, nil)
			if err != nil {
				t.Fatal(err)
			}

			mu.Lock()
			conn = c
			mu.Unlock()
		})

		srv.ListenAndServe()
	}()

	httpTsp, err := transport.Create("ws://"+addr, &transport.Params{
		ReconnectionInterval: 500 * time.Millisecond,
		ReconnectionTries:    1000000,
	})
	if err != nil {
		t.Fatal(err)
	}
	transportadapter.New(httpTsp)

	mu.Lock()
	conn.Close()
	mu.Unlock()

	srv.Close()

	done := make(chan struct{})
	go func() {
		http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			close(done)
		}))
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Reconnection timeout")
	}
}

func TestWSResubscription(t *testing.T) {
	var srv *http.Server
	var conn *websocket.Conn
	port := rand.Int31n(45535) + 20000
	addr := fmt.Sprintf("localhost:%d", port)

	go func() {
		srv = &http.Server{
			Addr: addr,
		}

		srv.Handler = testWSAuthSubscribeHandler(t, func(c *websocket.Conn) {
			conn = c
		})

		srv.ListenAndServe()
	}()

	httpTsp, err := transport.Create("ws://"+addr, &transport.Params{
		ReconnectionInterval: 500 * time.Millisecond,
		ReconnectionTries:    1000000,
	})
	if err != nil {
		t.Fatal(err)
	}
	adapter := transportadapter.New(httpTsp)

	_, err = adapter.Authenticate("jwt.token.123", 0)
	if err != nil {
		t.Fatal(err)
	}

	s, _, tspErr := adapter.Subscribe("subscribeCommands", 0, nil)
	if tspErr != nil {
		t.Fatal(tspErr)
	}

	conn.Close()
	srv.Close()

	go func() {
		http.ListenAndServe(addr, testWSAuthSubscribeHandler(t, nil))
	}()

	select {
	case <-s.DataChan:
	case <-time.After(3 * time.Second):
		t.Fatal("Resubscription timeout")
	}
}

func testWSAuthSubscribeHandler(t *testing.T, onConnect func(c *websocket.Conn)) http.HandlerFunc {
	u := &websocket.Upgrader{}

	return func(w http.ResponseWriter, r *http.Request) {
		c, err := u.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal(err)
		}

		if onConnect != nil {
			onConnect(c)
		}

		req := make(map[string]interface{})

		err = c.ReadJSON(&req)
		if err != nil {
			t.Fatal(err)
		}

		err = c.WriteJSON(authStubResponse(req))
		if err != nil {
			t.Fatal(err)
		}

		err = c.ReadJSON(&req)
		if err != nil {
			t.Fatal(err)
		}

		err = c.WriteJSON(subscriptionStubResponse(req))
		if err != nil {
			t.Fatal(err)
		}

		err = c.WriteJSON(subscriptionEventStub())
		if err != nil {
			t.Fatal(err)
		}
	}
}

func authStubResponse(req map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":    req["action"],
		"requestId": req["requestId"],
		"status":    "success",
	}
}

func subscriptionStubResponse(req map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":         req["action"],
		"requestId":      req["requestId"],
		"status":         "success",
		"subscriptionId": 1,
	}
}

func subscriptionEventStub() map[string]interface{} {
	return map[string]interface{}{
		"action":         "command/subscribe",
		"subscriptionId": 1,
		"command": map[string]interface{}{
			"id":   1,
			"name": "test",
		},
	}
}
