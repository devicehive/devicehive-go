package dh_test

import (
	"flag"
	"github.com/devicehive/devicehive-go/dh"
	"os"
	"testing"
	"fmt"
	"time"
)

var serverAddr = flag.String("serverAddress", "ws://localhost/api/websocket", "Server address without trailing slash")
var accessToken = flag.String("accessToken", "", "Your access token")
var refreshToken = flag.String("refreshToken", "", "Your refresh token")
var userId = flag.Int("userId", 0, "DH user ID")

var client *dh.Client

var waitTimeout time.Duration

func TestMain(m *testing.M) {
	flag.Parse()

	if *accessToken == "" || *refreshToken == "" || *userId == 0 {
		fmt.Println("Please provide accessToken, refreshToken and userId")
		os.Exit(1)
	}

	var err *dh.Error
	client, err = dh.ConnectWithToken(*serverAddr, *accessToken, *refreshToken)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	client.PollingWaitTimeoutSeconds = 7

	waitTimeout = time.Duration(client.PollingWaitTimeoutSeconds+1) * time.Second

	res := m.Run()
	os.Exit(res)
}
