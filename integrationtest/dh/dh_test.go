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

var accessToken = flag.String("accessToken", "eyJhbGciOiJIUzI1NiJ9.eyJwYXlsb2FkIjp7ImEiOlsyLDMsNCw1LDYsNyw4LDksMTAsMTEsMTIsMTUsMTYsMTddLCJlIjoxNTI2MDUwNDM5NDkzLCJ0IjoxLCJ1IjozNzg3NiwibiI6WyI0MTY5MSJdLCJkdCI6WyIqIl19fQ.SN9LoD4WXz3_KPXlifhRSYNaFPu7RgAZMlTSDEdl3IY", "Your access token")
var refreshToken = flag.String("refreshToken", "eyJhbGciOiJIUzI1NiJ9.eyJwYXlsb2FkIjp7ImEiOlsyLDMsNCw1LDYsNyw4LDksMTAsMTEsMTIsMTUsMTYsMTddLCJlIjoxNTQxNzczNDM5NDk0LCJ0IjowLCJ1IjozNzg3NiwibiI6WyI0MTY5MSJdLCJkdCI6WyIqIl19fQ.qhOH0W1HYBLE8tWNKUy1VjKT3xWHXdDaOLXA4oX_m1k", "Your refresh token")
var userId = flag.Int("userId", 0, "DH user ID")

var client *dh.Client
var serverTimestamp time.Time

var waitTimeout time.Duration

func TestMain(m *testing.M) {
	flag.Parse()

	var err *dh.Error
	client, err = dh.ConnectWithToken(httpServerAddr, *accessToken, *refreshToken)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	info, err := client.GetInfo()

	if err != nil {
		fmt.Printf("Couldn't get server info: %s\n", err)
		panic(err)
	}

	serverTimestamp = info.ServerTimestamp.Time

	client.PollingWaitTimeoutSeconds = 7

	waitTimeout = time.Duration(client.PollingWaitTimeoutSeconds+10) * time.Second

	res := m.Run()
	os.Exit(res)
}
