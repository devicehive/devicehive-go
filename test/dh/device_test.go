package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/matryer/is"
	"testing"
	"time"
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
	data := map[string]interface{}{
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

func TestDeviceSave(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	client := connect(addr)

	device, err := client.GetDevice("device-id")
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	device.Name = "updated name"

	err = device.Save()
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}
}

func TestDeviceListCommands(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	client := connect(addr)

	is := is.New(t)

	device, err := client.GetDevice("device-id")
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	listReqParams := &dh.ListParams{
		Start:     time.Now().Add(-1 * time.Hour),
		End:       time.Now(),
		Command:   "test command",
		Status:    "created",
		SortField: "timestamp",
		SortOrder: "ASC",
		Take:      10,
		Skip:      5,
	}
	list, err := device.ListCommands(listReqParams)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(len(list) != 0)
}

func TestDeviceSendCommand(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client := connect(addr)

	device, err := client.GetDevice("device-id")
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	comm, err := device.SendCommand("command name", nil, 120, time.Now(), "created", nil)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.Equal(comm.DeviceId, "device-id")
	is.Equal(comm.Command, "command name")
	is.True(comm.Id != 0)
	is.True(comm.LastUpdated.Unix() > 0)
	is.True(comm.UserId != 0)
}

func TestDeviceListNotifications(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	client := connect(addr)

	is := is.New(t)

	device, err := client.GetDevice("device-id")
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	listParams := &dh.ListParams{
		Start:        time.Now().Add(-1 * time.Hour),
		End:          time.Now(),
		Notification: "test notif",
		SortField:    "timestamp",
		SortOrder:    "ASC",
		Take:         10,
		Skip:         5,
	}
	list, err := device.ListNotifications(listParams)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(len(list) > 0)
}

func TestDeviceSendNotification(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	client := connect(addr)

	is := is.New(t)

	device, err := client.GetDevice("device-id")
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	notif, err := device.SendNotification("test notif", nil, time.Now())
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(notif != nil)
}
