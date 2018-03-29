package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/devicehive/devicehive-go/test/utils"
	"github.com/gorilla/websocket"
	"github.com/matryer/is"
	"os"
	"testing"
)

const serverAddr = "localhost:7357"
const wsServerAddr = "ws://" + serverAddr

var client *dh.Client
var resStub = utils.ResponseStub

func TestMain(m *testing.M) {
	res := m.Run()
	os.Exit(res)
}

func TestAuthenticate(t *testing.T) {
	is := is.New(t)

	srv := utils.TestWSServer(serverAddr, func(conn *websocket.Conn) {
		req := make(map[string]string)
		err := conn.ReadJSON(&req)

		if err != nil {
			panic(err)
		}

		is.Equal(req["action"], "authenticate")
		is.True(req["requestId"] != "")
		is.Equal(req["token"], "someTestToken")

		err = conn.WriteJSON(resStub.Authenticate(req["requestId"]))

		if err != nil {
			panic(err)
		}
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	res, err := client.Authenticate("someTestToken")

	is.NoErr(err)
	is.True(res)
}

func TestConnectionClose(t *testing.T) {
	is := is.New(t)

	srv := utils.TestWSServer(serverAddr, func(conn *websocket.Conn) {
		conn.ReadMessage()
		panic(nil)
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	_, err = client.Authenticate("test")

	is.Equal(err.Error(), "connection closed")
}

func TestInvalidResponse(t *testing.T) {
	is := is.New(t)

	srv := utils.TestWSServer(serverAddr, func(conn *websocket.Conn) {
		conn.ReadMessage()
		conn.WriteMessage(websocket.TextMessage, []byte("invalid response"))
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	_, err = client.Authenticate("test")

	is.Equal(err.Error(), "invalid service response")
}

func TestTokenByCreds(t *testing.T) {
	is := is.New(t)

	srv := utils.TestWSServer(serverAddr, func(conn *websocket.Conn) {
		req := make(map[string]string)
		err := conn.ReadJSON(&req)

		if err != nil {
			panic(err)
		}

		is.Equal(req["action"], "token")
		is.True(req["requestId"] != "")
		is.Equal(req["login"], "dhadmin")
		is.Equal(req["password"], "dhadmin_#911")

		err = conn.WriteJSON(resStub.Token(req["requestId"], "accTok", "refTok"))

		if err != nil {
			panic(err)
		}
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	accessToken, refreshToken, err := client.TokenByCreds("dhadmin", "dhadmin_#911")

	is.NoErr(err)
	is.Equal(accessToken, "accTok")
	is.Equal(refreshToken, "refTok")
}

func TestTokenRefresh(t *testing.T) {
	is := is.New(t)

	srv := utils.TestWSServer(serverAddr, func(conn *websocket.Conn) {
		req := make(map[string]string)
		err := conn.ReadJSON(&req)

		if err != nil {
			panic(err)
		}

		is.Equal(req["action"], "token/refresh")
		is.True(req["requestId"] != "")
		is.Equal(req["refreshToken"], "test refresh token")

		err = conn.WriteJSON(resStub.TokenRefresh(req["requestId"], "accTok"))

		if err != nil {
			panic(err)
		}
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	accessToken, err := client.TokenRefresh("test refresh token")

	is.NoErr(err)
	is.Equal(accessToken, "accTok")
}
