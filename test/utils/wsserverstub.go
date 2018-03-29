package utils

import (
	"net/http/httptest"
	"net"
	"net/http"
	"github.com/gorilla/websocket"
)

func TestWSServer(addr string, wsHandler func(c *websocket.Conn) ) *httptest.Server {
	l, err := net.Listen("tcp", addr)

	if err != nil {
		panic(err)
	}

	if wsHandler == nil {
		wsHandler = defaultWSHandler
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := upgrade(w, r)
		wsHandler(c)
	})
	srv := httptest.NewUnstartedServer(h)
	srv.Listener = l
	srv.Start()

	return srv
}

func defaultWSHandler(c *websocket.Conn) {
	for {
		_, msg, err := c.ReadMessage()

		if err != nil {
			panic(err)
			return
		}

		err = c.WriteMessage(0, msg)

		if err != nil {
			panic(err)
			return
		}
	}
}

var wsUpgrader = websocket.Upgrader{}

func upgrade(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	c, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		panic(err)
	}

	return c
}