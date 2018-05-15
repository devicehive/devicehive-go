package dh_test

import (
	"flag"
	"github.com/devicehive/devicehive-go/dh"
	"os"
	"testing"
	"fmt"
	"time"
)

var serverAddr = flag.String("serverAddress", "ws://playground-dev.devicehive.com/api/websocket", "Server address without trailing slash")
var accessToken = flag.String("accessToken", "eyJhbGciOiJIUzI1NiJ9.eyJwYXlsb2FkIjp7ImEiOlsyLDMsNCw1LDYsNyw4LDksMTAsMTEsMTIsMTUsMTYsMTddLCJlIjoxNTI2NTA0NDAwMDAwLCJ0IjoxLCJ1IjozNzg3NiwibiI6WyI0MTY5MSJdLCJkdCI6WyIqIl19fQ.MMW7Cz83ihdYAaQ0d84XAzF9KvRvOwacVRpuxAzp8n8", "Your access token")
var refreshToken = flag.String("refreshToken", "eyJhbGciOiJIUzI1NiJ9.eyJwYXlsb2FkIjp7ImEiOlsyLDMsNCw1LDYsNyw4LDksMTAsMTEsMTIsMTUsMTYsMTddLCJlIjoxNTI2NTA0NDAwMDAwLCJ0IjowLCJ1IjozNzg3NiwibiI6WyI0MTY5MSJdLCJkdCI6WyIqIl19fQ.DqBcgK079dZH0bEU-LcXduRpB0dKn6ql59xZ5MPxVzY", "Your refresh token")
var userId = flag.Int("userId", 37876, "DH user ID")

var client *dh.Client

var waitTimeout time.Duration

func TestMain(m *testing.M) {
	flag.Parse()

	if *accessToken == "" || *refreshToken == "" || *userId == 0 {
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
