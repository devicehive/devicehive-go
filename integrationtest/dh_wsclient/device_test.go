package dh_wsclient_test

import (
	"encoding/json"
	"github.com/devicehive/devicehive-go"
	"github.com/matryer/is"
	"testing"
	"time"
)

const testDeviceId = "go-test-dev"

func TestWSClientDevice(t *testing.T) {
	is := is.New(t)

	err := wsclient.PutDevice(testDeviceId, "", nil, 0, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	defer func() {
		err = wsclient.DeleteDevice(testDeviceId)
		if err != nil {
			t.Fatal(err)
		}

		testResponse(t, nil)
	}()

	err = wsclient.UpdateDevice(testDeviceId, &devicehive_go.Device{Name: "go-test-dev-updated"})
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	err = wsclient.GetDevice(testDeviceId)
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

	err := wsclient.PutDevice(testDeviceId, "", nil, 0, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	defer func() {
		err = wsclient.DeleteDevice(testDeviceId)
		if err != nil {
			t.Fatal(err)
		}

		testResponse(t, nil)
	}()

	err = wsclient.SendDeviceCommand(testDeviceId, "go-test", nil, 120, time.Time{}, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	err = wsclient.ListDeviceCommands(testDeviceId, nil)
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

	err := wsclient.PutDevice(testDeviceId, "", nil, 0, 0, false)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	defer func() {
		err = wsclient.DeleteDevice(testDeviceId)
		if err != nil {
			t.Fatal(err)
		}

		testResponse(t, nil)
	}()

	err = wsclient.SendDeviceNotification(testDeviceId, "go-test", nil, time.Time{})
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	err = wsclient.ListDeviceNotifications(testDeviceId, nil)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, func(data []byte) {
		var list []*devicehive_go.Notification
		json.Unmarshal(data, &list)
		is.True(len(list) > 0)
	})
}
