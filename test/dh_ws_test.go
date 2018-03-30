package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/devicehive/devicehive-go/test/utils"
	"github.com/gorilla/websocket"
	"github.com/matryer/is"
	"os"
	"testing"
	"time"
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

func TestRequestId(t *testing.T) {
	is := is.New(t)
	srv := utils.TestWSServer(serverAddr, func(conn *websocket.Conn) {
		req := make(map[string]interface{})
		conn.ReadJSON(&req)

		switch req["requestId"].(type) {
		case string:
			is.True(req["requestId"] != "")
		default:
			t.Error("requestId is not a string")
			is.Fail()
		}

		conn.WriteMessage(websocket.TextMessage, []byte("dummy response"))
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	// @TODO Maybe other methods should be placed here as well
	client.TokenRefresh("refresh token")
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

func TestTokenByPayload(t *testing.T) {
	is := is.New(t)

	srv := utils.TestWSServer(serverAddr, func(conn *websocket.Conn) {
		req := make(map[string]interface{})
		err := conn.ReadJSON(&req)

		if err != nil {
			panic(err)
		}

		is.Equal(req["action"], "token/create")

		err = conn.WriteJSON(resStub.Token(req["requestId"].(string), "accTok", "refTok"))

		if err != nil {
			panic(err)
		}
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	userId := 123
	actions := []string{ "ManageToken", "ManageNetworks" }
	networkIds := []string{ "n1", "n2" }
	deviceTypeIds := []string{ "d1", "d2" }
	expiration := time.Now()
	accessToken, refreshToken, err := client.TokenByPayload(userId, actions, networkIds, deviceTypeIds, &expiration)

	is.NoErr(err)
	is.Equal(accessToken, "accTok")
	is.Equal(refreshToken, "refTok")
}

func TestErrorResponseTokenByPayload(t *testing.T) {
	is := is.New(t)

	srv := utils.TestWSServer(serverAddr, func(conn *websocket.Conn) {
		req := make(map[string]interface{})
		err := conn.ReadJSON(&req)

		if err != nil {
			panic(err)
		}

		err = conn.WriteJSON(resStub.Unauthorized(req["action"].(string), req["requestId"].(string)))

		if err != nil {
			panic(err)
		}
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	_, _, err = client.TokenByPayload(1, nil, nil, nil, nil)

	is.Equal(err.Error(), "401 unauthorized")
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
