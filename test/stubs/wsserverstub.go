package stubs

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
)

type wsReqHandler func(reqData map[string]interface{}, conn *websocket.Conn) map[string]interface{}

type WSTestServer struct {
	handler wsReqHandler
	srv     *httptest.Server
}

func (wss *WSTestServer) Start() string {
	if wss.handler == nil {
		wss.handler = defaultWSHandler
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := upgrade(w, r)

		req := make(map[string]interface{})
		err := c.ReadJSON(&req)

		if err != nil {
			panic(err)
		}

		res := wss.handler(req, c)

		if res != nil {
			err = c.WriteJSON(res)

			if err != nil {
				panic(err)
			}
		}
	})
	srv := httptest.NewServer(h)

	wss.srv = srv

	return strings.Replace(srv.URL, "http", "ws", 1)
}

func (wss *WSTestServer) Close() {
	wss.srv.Close()
}

func (wss *WSTestServer) SetHandler(h wsReqHandler) {
	wss.handler = h
}

func defaultWSHandler(reqData map[string]interface{}, conn *websocket.Conn) map[string]interface{} {
	return ResponseStub.Respond(reqData)
}

var wsUpgrader = websocket.Upgrader{}

func upgrade(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	c, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		panic(err)
	}

	return c
}
