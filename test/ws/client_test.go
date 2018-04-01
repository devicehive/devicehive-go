package dh_ws_test

import (
	"testing"
	"github.com/matryer/is"
	"github.com/devicehive/devicehive-go/test/utils"
	"github.com/devicehive/devicehive-go/dh"
	"github.com/gorilla/websocket"
)

func TestAuthenticate(t *testing.T) {
	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		is.Equal(reqData["action"], "authenticate")
		return utils.ResponseStub.Authenticate(reqData["requestId"].(string))
	})

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	res, dhErr := client.Authenticate("someTestToken")

	logDHErr(t, dhErr)

	is.True(res)
}

func TestConnectionClose(t *testing.T) {
	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		panic(nil)
	})

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	_, dhErr := client.Authenticate("test")

	is.Equal(dhErr.Name(), dh.ConnClosedErr)
}

func TestInvalidResponse(t *testing.T) {
	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		c.WriteMessage(websocket.TextMessage, []byte("invalid response"))
		return nil
	})

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	_, dhErr := client.Authenticate("test")

	is.Equal(dhErr.Name(), dh.InvalidResponseErr)
}