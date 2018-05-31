package dh_wsclient_test

import (
	"github.com/devicehive/devicehive-go"
	"testing"
)

func TestCreateDeviceType(t *testing.T) {
	err := wsclient.CreateDeviceType("Test_Network", "Network for tests")

	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)
}

func TestGetDeviceType(t *testing.T) {
	err := wsclient.GetDeviceType(111)
	if err != nil {
		t.Fatal(err)
	}
	testResponse(t, nil)
}

func TestListDeviceTypes(t *testing.T) {
	err := wsclient.ListDeviceTypes(&devicehive_go.ListParams{})
	if err != nil {
		t.Fatal(err)
	}
	testResponse(t, nil)
}
