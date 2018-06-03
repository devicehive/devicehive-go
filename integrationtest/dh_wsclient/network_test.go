package dh_wsclient_test

import (
	"encoding/json"
	"github.com/devicehive/devicehive-go"
	"testing"
)

func TestCreateNetwork(t *testing.T) {
	err := wsclient.CreateNetwork("Test_Network", "Network for tests")
	if err != nil {
		t.Fatal(err)
	}

	networkId := 0
	testResponse(t, func(data []byte) {
		res := make(map[string]int)
		json.Unmarshal(data, &res)
		networkId = res["id"]
	})
	defer func() {
		err = wsclient.DeleteNetwork(networkId)
		if err != nil {
			t.Fatal(err)
		}
		testResponse(t, nil)
	}()

	err = wsclient.GetNetwork(networkId)
	if err != nil {
		t.Fatal(err)
	}
	testResponse(t, nil)

	err = wsclient.ListNetworks(&devicehive_go.ListParams{})
	if err != nil {
		t.Fatal(err)
	}
	testResponse(t, nil)
}
