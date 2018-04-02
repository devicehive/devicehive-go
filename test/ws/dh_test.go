package dh_ws_test

import (
	"github.com/devicehive/devicehive-go/test/stubs"
	"os"
	"testing"
)

const serverAddr = "localhost:7357"
const wsServerAddr = "ws://" + serverAddr

var wsTestSrv = &stubs.WSTestServer{}

func TestMain(m *testing.M) {
	wsTestSrv.Start(serverAddr)
	defer wsTestSrv.Close()

	res := m.Run()
	os.Exit(res)
}
