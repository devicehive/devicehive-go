package dh_test

import (
	"flag"
	"github.com/devicehive/devicehive-go/dh"
	"os"
	"testing"
)

const serverAddr = "playground-dev.devicehive.com/api/websocket"
const wsServerAddr = "ws://" + serverAddr

var dhLogin = flag.String("dhLogin", "dhadmin", "Your username")
var dhPass = flag.String("dhPassword", "dhadmin_#911", "Your password")

var client *dh.Client

func TestMain(m *testing.M) {
	flag.Parse()

	var err *dh.Error
	client, err = dh.ConnectWithCreds(wsServerAddr, *dhLogin, *dhPass)

	if err != nil {
		panic(err)
	}

	res := m.Run()
	os.Exit(res)
}
