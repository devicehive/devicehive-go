package transportadapter

import (
	"testing"
	"net/http"
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/internal/transportadapter"
	"github.com/gorilla/websocket"
	"time"
)

func TestWSAdapterReconnection(t *testing.T) {
	var srv *http.Server
	var conn *websocket.Conn
	go func() {
		srv = &http.Server{
			Addr: "localhost:1337",
		}

		u := &websocket.Upgrader{}
		srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := u.Upgrade(w, r, nil)
			if err != nil {
				t.Fatal(err)
			}

			conn = c
		})

		srv.ListenAndServe()
	}()

	addr := "ws://localhost:1337"
	httpTsp, err := transport.Create(addr, &transport.Params{
		ReconnectionInterval: 500 * time.Millisecond,
		ReconnectionTries: 1000000,
	})
	if err != nil {
		t.Fatal(err)
	}
	transportadapter.New(httpTsp)

	conn.Close()
	srv.Close()

	done := make(chan struct{})
	go func() {
		http.ListenAndServe("localhost:1337", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			close(done)
		}))
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Reconnection timeout")
	}
}
