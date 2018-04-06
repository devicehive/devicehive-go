package dh_test

import (
	"github.com/devicehive/devicehive-go/test/stubs"
	"os"
	"testing"
	"github.com/gorilla/websocket"
)

const serverAddr = "localhost:7357"
const wsServerAddr = "ws://" + serverAddr

var wsTestSrv = &stubs.WSTestServer{}

func TestMain(m *testing.M) {
	wsTestSrv.Start(serverAddr)
	defer wsTestSrv.Close()

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, conn *websocket.Conn) map[string]interface{} {
		return stubs.ResponseStub.Respond(reqData)
	})

	res := m.Run()
	os.Exit(res)
}
