package dh_ws_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/devicehive/devicehive-go/testutils"
	"github.com/gorilla/websocket"
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestTokenByCreds(t *testing.T) {
	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		is.Equal(reqData["action"], "token")
		return stubs.ResponseStub.Token(reqData["requestId"].(string), "accTok", "refTok")
	})

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	accessToken, refreshToken, dhErr := client.TokenByCreds("dhadmin", "dhadmin_#911")

	testutils.LogDHErr(t, dhErr)

	is.Equal(accessToken, "accTok")
	is.Equal(refreshToken, "refTok")
}

func TestTokenByPayload(t *testing.T) {
	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		is.Equal(reqData["action"], "token/create")
		return stubs.ResponseStub.Token(reqData["requestId"].(string), "accTok", "refTok")
	})

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

	testutils.LogDHErr(t, dhErr)

	is.Equal(accessToken, "accTok")
	is.Equal(refreshToken, "refTok")
}

func TestErrorResponseTokenByPayload(t *testing.T) {
	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		return stubs.ResponseStub.Unauthorized(reqData["action"].(string), reqData["requestId"].(string))
	})

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

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		is.Equal(reqData["action"], "token/refresh")
		return stubs.ResponseStub.TokenRefresh(reqData["requestId"].(string), "accTok")
	})

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	accessToken, dhErr := client.TokenRefresh("test refresh token")

	testutils.LogDHErr(t, dhErr)

	is.Equal(accessToken, "accTok")
}