package dh_wsclient_test

import (
	"encoding/json"
	"github.com/devicehive/devicehive-go"
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestWSClientDevice(t *testing.T) {
	is := is.New(t)

	err := wsclient.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	defer func() {
		err = wsclient.DeleteDevice("go-test-dev")
		if err != nil {
			t.Fatal(err)
		}

		testResponse(t, nil)
	}()

	err = wsclient.UpdateDevice("go-test-dev", &devicehive_go.Device{Name: "go-test-dev-updated"})
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	err = wsclient.GetDevice("go-test-dev")
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, func(data []byte) {
		dev := &devicehive_go.Device{}
		json.Unmarshal(data, dev)
		is.Equal(dev.Name, "go-test-dev-updated")
	})

	err = wsclient.ListDevices(nil)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, func(data []byte) {
		var list []*devicehive_go.Device
		json.Unmarshal(data, &list)
		is.True(len(list) > 0)
	})
}

func TestWSClientDeviceCommands(t *testing.T) {
	is := is.New(t)

	err := wsclient.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	defer func() {
		err = wsclient.DeleteDevice("go-test-dev")
		if err != nil {
			t.Fatal(err)
		}

		testResponse(t, nil)
	}()

	err = wsclient.SendDeviceCommand("go-test-dev", "go-test", nil, 120, time.Time{}, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	err = wsclient.ListDeviceCommands("go-test-dev", nil)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, func(data []byte) {
		var list []*devicehive_go.Command
		json.Unmarshal(data, &list)
		is.True(len(list) > 0)
	})
}

func TestWSClientDeviceNotifications(t *testing.T) {
	is := is.New(t)

	err := wsclient.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	defer func() {
		err = wsclient.DeleteDevice("go-test-dev")
		if err != nil {
			t.Fatal(err)
		}

		testResponse(t, nil)
	}()

	err = wsclient.SendDeviceNotification("go-test-dev", "go-test", nil, time.Time{})
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	err = wsclient.ListDeviceNotifications("go-test-dev", nil)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, func(data []byte) {
		var list []*devicehive_go.Notification
		json.Unmarshal(data, &list)
		is.True(len(list) > 0)
	})
}
