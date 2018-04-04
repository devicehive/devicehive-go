package internal_test

import (
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/gorilla/websocket"
	"github.com/matryer/is"
	"testing"
	"time"
)

const serverAddr = "localhost:7358"
const wsServerAddr = "ws://" + serverAddr
const testTimeout = 300 * time.Millisecond

func TestRequestId(t *testing.T) {
	is := is.New(t)
	wsTestSrv, srvClose := createWStestSrv(serverAddr)
	defer srvClose()

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		is.True(reqData["requestId"] != "")
		return nil
	})

	wsTsp, err := transport.Create(wsServerAddr)

	is.NoErr(err)

	wsTsp.Request(map[string]interface{}{}, testTimeout)
}

func TestTimeout(t *testing.T) {
	is := is.New(t)
	wsTestSrv, srvClose := createWStestSrv(serverAddr)
	defer srvClose()

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		<-time.After(testTimeout + 1*time.Second)

		return map[string]interface{}{
			"result": "success",
		}
	})

	wsTsp, err := transport.Create(wsServerAddr)

	is.NoErr(err)

	res, tspErr := wsTsp.Request(map[string]interface{}{}, testTimeout)

	is.True(res == nil)
	is.Equal(tspErr.Name(), transport.TimeoutErr)
}

func TestInvalidResponse(t *testing.T) {
	is := is.New(t)
	wsTestSrv, srvClose := createWStestSrv(serverAddr)
	defer srvClose()

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		c.WriteMessage(websocket.TextMessage, []byte("invalid response"))
		return nil
	})

	wsTsp, err := transport.Create(wsServerAddr)

	is.NoErr(err)

	res, tspErr := wsTsp.Request(map[string]interface{}{}, testTimeout)

	is.True(res == nil)
	is.Equal(tspErr.Name(), transport.TimeoutErr)
}

func TestConnectionClose(t *testing.T) {
	is := is.New(t)
	wsTestSrv, srvClose := createWStestSrv(serverAddr)
	defer srvClose()

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		c.Close()
		return nil
	})

	wsTsp, err := transport.Create(wsServerAddr)

	is.NoErr(err)

	res, tspErr := wsTsp.Request(map[string]interface{}{}, testTimeout)

	is.True(res == nil)
	is.Equal(tspErr.Name(), transport.ConnClosedErr)
}

func createWStestSrv(serverAddr string) (srv *stubs.WSTestServer, srvClose func()) {
	wsTestSrv := &stubs.WSTestServer{}

	wsTestSrv.Start(serverAddr)
	return wsTestSrv, wsTestSrv.Close
}
