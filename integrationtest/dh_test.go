package integrationtest

import (
	"testing"
	"github.com/matryer/is"
	"github.com/devicehive/devicehive-go/dh"
	"os"
	"fmt"
	"flag"
)

const serverAddr = "playground-dev.devicehive.com/api/websocket"
const wsServerAddr = "ws://" + serverAddr

var tok = flag.String("accessToken", "", "Your DeviceHive access token")

func TestMain(m *testing.M) {
	flag.Parse()
	if *tok == "" {
		fmt.Println("Access token is not specified")
		os.Exit(1)
	}

	res := m.Run()
	os.Exit(res)
}

func TestWSConnection(t *testing.T) {
	is := is.New(t)

	client, err := dh.Connect(wsServerAddr)

	is.NoErr(err)
	is.True(client != nil)
}

func TestAuthenticate(t *testing.T) {
	is := is.New(t)

	client, err := dh.Connect(wsServerAddr)

	is.NoErr(err)

	res, err := client.Authenticate(*tok)

	is.NoErr(err)
	is.True(res)
}