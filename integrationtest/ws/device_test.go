package dh_test

import (
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestDevice(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
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

	device, err = client.GetDevice(device.Id)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(device != nil)
	is.Equal(device.Name, "updated name")

	err = device.Remove()
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}
}

func TestDeviceCommands(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	comm, err := device.SendCommand("test command", nil, 5, time.Time{}, "", nil)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(comm != nil)

	comm.Status = "updated"

	err = comm.Save()
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	list, err := device.ListCommands(nil)

	is.True(len(list) > 0)
	is.Equal(list[0].Status, "updated")

	err = device.Remove()
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}
}

func TestDeviceNotifications(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	notif, err := device.SendNotification("test notif", nil, time.Time{})
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(notif != nil)

	list, err := device.ListNotifications(nil)

	is.True(len(list) > 0)

	err = device.Remove()
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}
}
