package utils

import (
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"net/http/httptest"
)

type wsReqHandler func(reqData map[string]interface{}, conn *websocket.Conn) map[string]interface{}

func TestWSServer(addr string, handler wsReqHandler) *httptest.Server {
	l, err := net.Listen("tcp", addr)

	if err != nil {
		panic(err)
	}

	if handler == nil {
		handler = defaultWSHandler
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := upgrade(w, r)
		defer c.Close()

		req := make(map[string]interface{})
		err := c.ReadJSON(&req)

		if err != nil {
			panic(err)
		}

		res := handler(req, c)

		if res != nil  {
			err = c.WriteJSON(res)

			if err != nil {
				panic(err)
			}
		}
	})
	srv := httptest.NewUnstartedServer(h)
	srv.Listener = l
	srv.Start()

	return srv
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
