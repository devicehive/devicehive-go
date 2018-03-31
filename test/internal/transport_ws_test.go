package internal_test

import (
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/test/utils"
	"github.com/gorilla/websocket"
	"github.com/matryer/is"
	"testing"
)

func TestRequestId(t *testing.T) {
	is := is.New(t)
	wsTestSrv := &utils.WSTestServer{}

	wsTestSrv.Start("localhost:7357")
	defer wsTestSrv.Close()

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		is.True(reqData["requestId"] != "")
		return nil
	})

	wsTsp, err := transport.Create("ws://localhost:7357")

	is.NoErr(err)

	wsTsp.Request(map[string]interface{}{})
}
