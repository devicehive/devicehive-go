package dh_wsclient_test

import (
	"errors"
	"flag"
	"fmt"
	"github.com/devicehive/devicehive-go"
	"os"
	"testing"
	"time"
)

var serverAddr = flag.String("serverAddress", "", "Server address without trailing slash")
var accessToken = flag.String("accessToken", "", "Your access token")
var refreshToken = flag.String("refreshToken", "", "Your refresh token")
var userId = flag.Int("userId", 0, "DH user ID")

var wsclient *devicehive_go.WSClient

const TestTimeout = 3 * time.Second

func TestMain(m *testing.M) {
	flag.Parse()

	if *accessToken == "" || *refreshToken == "" || *userId == 0 {
		os.Exit(1)
	}

	var err *devicehive_go.Error
	wsclient, err = devicehive_go.WSConnect(*serverAddr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = wsclient.Authenticate(*accessToken)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	select {
	case <-wsclient.DataChan:
	case err := <-wsclient.ErrorChan:
		fmt.Println(err)
		panic(err)
	case <-time.After(TestTimeout):
		fmt.Println("Timeout")
		panic(errors.New("timeout"))
	}

	res := m.Run()
	os.Exit(res)
}

func testResponse(t *testing.T, dataCallback func([]byte)) {
	select {
	case data := <-wsclient.DataChan:
		if dataCallback != nil {
			dataCallback(data)
		}
	case err := <-wsclient.ErrorChan:
		t.Fatal(err)
	case <-time.After(TestTimeout):
		t.Fatal("Timeout")
	}
}
