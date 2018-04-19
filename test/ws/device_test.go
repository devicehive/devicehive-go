package dh_test

import (
	"testing"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/matryer/is"
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

	name := "device-name"
	data := map[string]interface{} {
		"param": "test",
	}
	networkId := int64(1)
	deviceTypeId := int64(1)
	isBlocked := false
	device, err := client.PutDevice("device-id", name, data, networkId, deviceTypeId, isBlocked)
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

	device, err := client.GetDevice("device-id")
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	err = device.Remove()
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}
}
