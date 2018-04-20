package dh_test

import (
	"testing"
	"github.com/matryer/is"
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

	comm, err := device.SendCommand("test command", nil, 0, time.Time{}, "", nil)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(comm != nil)

	list, err := device.ListCommands(nil)

	is.True(len(list) > 0)

	err = device.Remove()
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}
}
