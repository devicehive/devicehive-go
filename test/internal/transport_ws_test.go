package internal_test

import (
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/gorilla/websocket"
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestRequestId(t *testing.T) {
	is := is.New(t)
	wsTestSrv := &stubs.WSTestServer{}

	wsTestSrv.Start("localhost:7357")
	defer wsTestSrv.Close()

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		is.True(reqData["requestId"] != "")
		return nil
	})

	wsTsp, err := transport.Create("ws://localhost:7357")

	is.NoErr(err)

	wsTsp.Request(map[string]interface{}{}, 0)
}

func TestTimeout(t *testing.T) {
	is := is.New(t)
	wsTestSrv := &stubs.WSTestServer{}

	wsTestSrv.Start("localhost:7357")
	defer wsTestSrv.Close()

	timeout := 300 * time.Millisecond

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		<-time.After(timeout + 1 * time.Second)

		return map[string]interface{} {
			"result": "success",
		}
	})

	wsTsp, err := transport.Create("ws://localhost:7357")

	is.NoErr(err)

	res, tspErr := wsTsp.Request(map[string]interface{}{}, timeout)

	is.True(res == nil)
	is.Equal(tspErr.Name(), transport.TimeoutErr)
}
