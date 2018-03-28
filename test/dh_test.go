package dh_test

import (
	"testing"
	"github.com/matryer/is"
	"github.com/devicehive/devicehive-go/dh"
	"os"
	"github.com/gorilla/websocket"
	"github.com/devicehive/devicehive-go/test/utils"
)

const serverAddr = "localhost:7357"
const wsServerAddr = "ws://" + serverAddr

func TestMain(m *testing.M) {
	res := m.Run()
	os.Exit(res)
}

func TestWSConnection(t *testing.T) {
	is := is.New(t)

	srv := utils.TestWSServer(serverAddr, nil)
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	is.NoErr(err)
	is.True(client != nil)
}

func TestAuthenticate(t *testing.T) {
	is := is.New(t)

	srv := utils.TestWSServer(serverAddr, func(conn *websocket.Conn) {
		req := make(map[string]string)
		err := conn.ReadJSON(&req)

		is.NoErr(err)
		is.Equal(req["action"], "authenticate")
		is.True(req["requestId"] != "")
		is.Equal(req["token"], "someTestToken")

		conn.WriteJSON(map[string]string{
			"action": req["action"],
			"requestId": req["requestId"],
			"status": "success",
		})
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	is.NoErr(err)

	res, err := client.Authenticate("someTestToken")

	is.NoErr(err)
	is.True(res)
}