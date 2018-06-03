package dh_wsclient_test

import (
	"encoding/json"
	"github.com/devicehive/devicehive-go"
	"testing"
)

func TestCreateDeviceType(t *testing.T) {
	err := wsclient.CreateDeviceType("Test_DeviceType", "Device type for tests")
	if err != nil {
		t.Fatal(err)
	}
	devTypeId := 0
	testResponse(t, func(data []byte) {
		res := make(map[string]int)
		json.Unmarshal(data, &res)
		devTypeId = res["id"]
	})
	defer func() {
		err = wsclient.DeleteDeviceType(devTypeId)
		if err != nil {
			t.Fatal(err)
		}
		testResponse(t, nil)
	}()

	err = wsclient.GetDeviceType(1)
	if err != nil {
		t.Fatal(err)
	}
	testResponse(t, nil)

	err = wsclient.ListDeviceTypes(&devicehive_go.ListParams{})
	if err != nil {
		t.Fatal(err)
	}
	testResponse(t, nil)
}
