package dh_ws_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/devicehive/devicehive-go/test/utils"
	"os"
	"testing"
)

const serverAddr = "localhost:7357"
const wsServerAddr = "ws://" + serverAddr

var wsTestSrv = &utils.WSTestServer{}

func TestMain(m *testing.M) {
	wsTestSrv.Start(serverAddr)
	defer wsTestSrv.Close()

	res := m.Run()
	os.Exit(res)
}

func logDHErr(t *testing.T, err *dh.Error) {
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}
}
