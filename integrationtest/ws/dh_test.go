package dh_ws_test

import (
	"flag"
	"github.com/devicehive/devicehive-go/dh"
	"os"
	"testing"
)

const serverAddr = "playground-dev.devicehive.com/api/websocket"
const wsServerAddr = "ws://" + serverAddr

var tok = flag.String("accessToken", "", "Your JWT access token")
var dhLogin = flag.String("dhLogin", "dhadmin", "Your username")
var dhPass = flag.String("dhPassword", "dhadmin_#911", "Your password")

var client *dh.Client

func TestMain(m *testing.M) {
	var err error
	client, err = dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	flag.Parse()

	res := m.Run()
	os.Exit(res)
}
