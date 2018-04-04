package stubs

import (
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
)

var mu = sync.Mutex{}

type wsReqHandler func(reqData map[string]interface{}, conn *websocket.Conn) map[string]interface{}

type WSTestServer struct {
	handler wsReqHandler
	srv     *httptest.Server
}

func (wss *WSTestServer) Start(addr string) {
	l, err := net.Listen("tcp", addr)

	if err != nil {
		panic(err)
	}

	if wss.handler == nil {
		wss.handler = defaultWSHandler
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := upgrade(w, r)
		defer c.Close()

		req := make(map[string]interface{})
		err := c.ReadJSON(&req)

		if err != nil {
			panic(err)
		}

		mu.Lock()
		defer mu.Unlock()
		res := wss.handler(req, c)

		if res != nil {
			err = c.WriteJSON(res)

			if err != nil {
				panic(err)
			}
		}
	})
	srv := httptest.NewUnstartedServer(h)
	srv.Listener = l
	srv.Start()

	wss.srv = srv
}

func (wss *WSTestServer) Close() {
	wss.srv.Close()
}

func (wss *WSTestServer) SetHandler(h wsReqHandler) {
	mu.Lock()
	defer mu.Unlock()
	wss.handler = h
}

func defaultWSHandler(req map[string]interface{}, conn *websocket.Conn) map[string]interface{} {
	return req
}

var wsUpgrader = websocket.Upgrader{}

func upgrade(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	c, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		panic(err)
	}

	return c
}
