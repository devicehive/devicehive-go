package stubs

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
	"log"
)

func StartWSTestServer() (srv *WSTestServer, addr string, closeSrv func()) {
	srv = &WSTestServer{}
	addr = srv.Start()

	return srv, addr, srv.Close
}

type wsReqHandler func(reqData map[string]interface{}, conn *websocket.Conn)

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

		for {
			req := make(map[string]interface{})
			err := c.ReadJSON(&req)

			if err != nil {
				log.Println("peer closed connection")
				return
			}

			wss.handler(req, c)
		}
	})
	srv := httptest.NewServer(h)

	wss.srv = srv

	return strings.Replace(srv.URL, "http", "ws", 1)
}

func (wss *WSTestServer) Close() {
	wss.srv.Close()
}

func (wss *WSTestServer) SetRequestHandler(h wsReqHandler) {
	wss.handler = h
}

func defaultWSHandler(reqData map[string]interface{}, conn *websocket.Conn) {
	err := conn.WriteJSON(ResponseStub.Respond(reqData))

	if err != nil {
		panic(err)
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
