package dh_test

import (
	"flag"
	"github.com/devicehive/devicehive-go/dh"
	"os"
	"testing"
)

const serverAddr = "playground-dev.devicehive.com/api"
const wsServerAddr = "ws://" + serverAddr + "/websocket"
const httpServerAddr = "http://" + serverAddr + "/rest"

var accessToken = flag.String("accessToken", "eyJhbGciOiJIUzI1NiJ9.eyJwYXlsb2FkIjp7ImEiOlsyLDMsNCw1LDYsNyw4LDksMTAsMTEsMTIsMTUsMTYsMTddLCJlIjoxNTI1OTQwOTcxNDI4LCJ0IjoxLCJ1IjozNzg3NiwibiI6WyI0MTY5MSJdLCJkdCI6WyIqIl19fQ.L6qyLIB_zuYnpQtrgWXgop2518eqXxJ26QuWyPeyhbA", "Your access token")
var refreshToken = flag.String("refreshToken", "eyJhbGciOiJIUzI1NiJ9.eyJwYXlsb2FkIjp7ImEiOlsyLDMsNCw1LDYsNyw4LDksMTAsMTEsMTIsMTUsMTYsMTddLCJlIjoxNTQxNjYzOTcxNDI4LCJ0IjowLCJ1IjozNzg3NiwibiI6WyI0MTY5MSJdLCJkdCI6WyIqIl19fQ.pyj5jqFGDiIMSKeAQuokc6X40DOSFuncxnU2bSuBo34", "Your refresh token")

var client *dh.Client

func TestMain(m *testing.M) {
	flag.Parse()

	var err *dh.Error
	client, err = dh.ConnectWithToken(wsServerAddr, *accessToken, *refreshToken)

	if err != nil {
		panic(err)
	}

	client.PollingWaitTimeoutSeconds = 7

	res := m.Run()
	os.Exit(res)
}
