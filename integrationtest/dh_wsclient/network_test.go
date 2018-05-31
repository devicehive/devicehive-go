package dh_wsclient_test

import (
	"github.com/devicehive/devicehive-go"
	"testing"
)

func TestCreateNetwork(t *testing.T) {
	err := wsclient.CreateNetwork("Test_Network", "Network for tests")

	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)
}

func TestGetNetwork(t *testing.T) {
	err := wsclient.GetNetwork(41691)
	if err != nil {
		t.Fatal(err)
	}
	testResponse(t, nil)
}

func TestListNetworks(t *testing.T) {
	err := wsclient.ListNetworks(&devicehive_go.ListParams{})
	if err != nil {
		t.Fatal(err)
	}
	testResponse(t, nil)
}
