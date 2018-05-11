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

var accessToken = flag.String("accessToken", "eyJhbGciOiJIUzI1NiJ9.eyJwYXlsb2FkIjp7ImEiOlsyLDMsNCw1LDYsNyw4LDksMTAsMTEsMTIsMTUsMTYsMTddLCJlIjoxNTI2MDM0OTU5NzA5LCJ0IjoxLCJ1IjozNzg3NiwibiI6WyI0MTY5MSJdLCJkdCI6WyIqIl19fQ.xZsDvbFrtW9WMjuYbH2CcROUuI5HJa-Zmxfu2RQ9T0w", "Your access token")
var refreshToken = flag.String("refreshToken", "eyJhbGciOiJIUzI1NiJ9.eyJwYXlsb2FkIjp7ImEiOlsyLDMsNCw1LDYsNyw4LDksMTAsMTEsMTIsMTUsMTYsMTddLCJlIjoxNTQxNzU3OTU5NzEwLCJ0IjowLCJ1IjozNzg3NiwibiI6WyI0MTY5MSJdLCJkdCI6WyIqIl19fQ.a8-zYF2yq7fy9YBKzx2qfCYH3bdHr4tp_RZjKtzOEEI", "Your refresh token")
var userId = flag.Int("userId", 0, "DH user ID")

var client *dh.Client
var serverTimestamp time.Time

func TestMain(m *testing.M) {
	flag.Parse()

	var err *dh.Error
	client, err = dh.ConnectWithToken(wsServerAddr, *accessToken, *refreshToken)

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

	res := m.Run()
	os.Exit(res)
}
