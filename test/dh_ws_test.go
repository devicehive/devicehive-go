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

	srv := utils.TestWSServer(serverAddr, func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		is.Equal(reqData["action"], "authenticate")
		return resStub.Authenticate(reqData["requestId"].(string))
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	res, dhErr := client.Authenticate("someTestToken")

	is.True(dhErr == nil)
	is.True(res)
}

func TestConnectionClose(t *testing.T) {
	is := is.New(t)

	srv := utils.TestWSServer(serverAddr, func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		panic(nil)
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	_, dhErr := client.Authenticate("test")

	is.Equal(dhErr.Name(), dh.ConnClosedErr)
}

func TestInvalidResponse(t *testing.T) {
	is := is.New(t)

	srv := utils.TestWSServer(serverAddr, func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		c.WriteMessage(websocket.TextMessage, []byte("invalid response"))
		return nil
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	_, dhErr := client.Authenticate("test")

	is.Equal(dhErr.Name(), dh.InvalidResponseErr)
}

func TestRequestId(t *testing.T) {
	is := is.New(t)
	srv := utils.TestWSServer(serverAddr, func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		switch reqData["requestId"].(type) {
		case string:
			is.True(reqData["requestId"] != "")
		default:
			t.Error("requestId is not a string")
			is.Fail()
		}

		c.WriteMessage(websocket.TextMessage, []byte("dummy response"))

		return nil
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

	srv := utils.TestWSServer(serverAddr, func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		is.Equal(reqData["action"], "token")
		return resStub.Token(reqData["requestId"].(string), "accTok", "refTok")
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	accessToken, refreshToken, dhErr := client.TokenByCreds("dhadmin", "dhadmin_#911")

	is.True(dhErr == nil)
	is.Equal(accessToken, "accTok")
	is.Equal(refreshToken, "refTok")
}

func TestTokenByPayload(t *testing.T) {
	is := is.New(t)

	srv := utils.TestWSServer(serverAddr, func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		is.Equal(reqData["action"], "token/create")
		return resStub.Token(reqData["requestId"].(string), "accTok", "refTok")
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	userId := 123
	actions := []string{"ManageToken", "ManageNetworks"}
	networkIds := []string{"n1", "n2"}
	deviceTypeIds := []string{"d1", "d2"}
	expiration := time.Now()
	accessToken, refreshToken, dhErr := client.TokenByPayload(userId, actions, networkIds, deviceTypeIds, &expiration)

	is.True(dhErr == nil)
	is.Equal(accessToken, "accTok")
	is.Equal(refreshToken, "refTok")
}

func TestErrorResponseTokenByPayload(t *testing.T) {
	is := is.New(t)

	srv := utils.TestWSServer(serverAddr, func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		return resStub.Unauthorized(reqData["action"].(string), reqData["requestId"].(string))
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	_, _, dhErr := client.TokenByPayload(1, nil, nil, nil, nil)

	is.Equal(dhErr.Name(), dh.ServiceErr)
	is.Equal(dhErr.Error(), "401 unauthorized")
}

func TestTokenRefresh(t *testing.T) {
	is := is.New(t)

	srv := utils.TestWSServer(serverAddr, func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		is.Equal(reqData["action"], "token/refresh")
		return resStub.TokenRefresh(reqData["requestId"].(string), "accTok")
	})
	defer srv.Close()

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	accessToken, dhErr := client.TokenRefresh("test refresh token")

	is.True(dhErr == nil)
	is.Equal(accessToken, "accTok")
}
