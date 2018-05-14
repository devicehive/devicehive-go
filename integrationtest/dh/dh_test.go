package dh_test

import (
	"flag"
	"github.com/devicehive/devicehive-go/dh"
	"os"
	"testing"
	"fmt"
	"time"
)

const serverAddr = "playground-dev.devicehive.com/api"
const wsServerAddr = "ws://" + serverAddr + "/websocket"
const httpServerAddr = "http://" + serverAddr + "/rest"

var accessToken = flag.String("accessToken", "", "Your access token")
var refreshToken = flag.String("refreshToken", "", "Your refresh token")
var userId = flag.Int("userId", 0, "DH user ID")

var client *dh.Client

var waitTimeout time.Duration

func TestMain(m *testing.M) {
	flag.Parse()

	var err *dh.Error
	client, err = dh.ConnectWithToken(wsServerAddr, *accessToken, *refreshToken)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	client.PollingWaitTimeoutSeconds = 7

	waitTimeout = time.Duration(client.PollingWaitTimeoutSeconds+10) * time.Second

	res := m.Run()
	os.Exit(res)
}
