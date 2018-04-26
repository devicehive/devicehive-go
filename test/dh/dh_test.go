package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	res := m.Run()
	os.Exit(res)
}

func connect(addr string) *dh.Client {
	client, err := dh.ConnectWithToken(addr, "", "")

	if err != nil {
		panic(err)
	}

	return client
}
