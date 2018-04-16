package dh_test

import (
	"flag"
	"github.com/devicehive/devicehive-go/dh"
	"os"
	"testing"
)

const serverAddr = "playground-dev.devicehive.com/api/websocket"
const wsServerAddr = "ws://" + serverAddr
const testDeviceId = "4NemW3PE9BHRSqb0DVVgsphZh7SCZzgm3Lxg"

var dhLogin = flag.String("dhLogin", "dhadmin", "Your username")
var dhPass = flag.String("dhPassword", "dhadmin_#911", "Your password")

var client *dh.Client

func TestMain(m *testing.M) {
	var err *dh.Error
	client, err = dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	flag.Parse()

	res := m.Run()
	os.Exit(res)
}

func auth() *dh.Error {
	accTok, _, err := client.TokenByCreds(*dhLogin, *dhPass)

	if err != nil {
		return err
	}

	_, err = client.Authenticate(accTok)

	if err != nil {
		return err
	}

	return nil
}
