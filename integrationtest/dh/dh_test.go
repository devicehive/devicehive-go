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

var accessToken = flag.String("accessToken", "eyJhbGciOiJIUzI1NiJ9.eyJwYXlsb2FkIjp7ImEiOlsyLDMsNCw1LDYsNyw4LDksMTAsMTEsMTIsMTUsMTYsMTddLCJlIjoxNTI1OTQ4MzYyNzczLCJ0IjoxLCJ1IjozNzg3NiwibiI6WyI0MTY5MSJdLCJkdCI6WyIqIl19fQ.Jn43WlrM-kHjYunF4Mm_V90LwmcbZ4AdyTDsui_PFF8", "Your access token")
var refreshToken = flag.String("refreshToken", "eyJhbGciOiJIUzI1NiJ9.eyJwYXlsb2FkIjp7ImEiOlsyLDMsNCw1LDYsNyw4LDksMTAsMTEsMTIsMTUsMTYsMTddLCJlIjoxNTQxNjcxMzYyNzczLCJ0IjowLCJ1IjozNzg3NiwibiI6WyI0MTY5MSJdLCJkdCI6WyIqIl19fQ.J56RR-LVivQBe8CSINloYjyhO3gLlWokKX1jnqRxQoc", "Your refresh token")
var userId = flag.Int("userId", 0, "DH user ID")

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
