package dh_test

import (
	"testing"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/matryer/is"
	"github.com/devicehive/devicehive-go/dh"
)

func TestGetDevice(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	client := connect(addr)

	is := is.New(t)

	device, err := client.GetDevice("device-id")
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(device != nil)
	is.True(device.Id != "")
	is.True(device.Name != "")
	is.True(device.NetworkId != 0)
	is.True(device.DeviceTypeId != 0)
}

func TestPutDevice(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	client := connect(addr)

	is := is.New(t)

	device := &dh.Device{
		Name: "device",
		Data: map[string]interface{} {
			"param": "test",
		},
		NetworkId: 1,
		DeviceTypeId: 1,
		IsBlocked: false,
	}
	err := client.PutDevice("device-id", device)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.Equal(device.Id, "device-id")
}

func TestDeviceRemove(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	client := connect(addr)

	err := client.RemoveDevice("device-id")
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}
}
