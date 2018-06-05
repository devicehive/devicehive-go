package dh_test

import (
	dh "github.com/devicehive/devicehive-go"
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestDevice(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	list, err := client.ListDevices(&dh.ListParams{
		NamePattern: "go-%-dev",
	})

	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	is.Equal(len(list), 1)

	device.Name = "updated name"
	err = device.Save()
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	device, err = client.GetDevice(device.Id)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	is.True(device != nil)
	is.Equal(device.Name, "updated name")

	err = device.Remove()
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}
}

func TestDeviceCommands(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-command", "", nil, 0, 0, false)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	comm, err := device.SendCommand("test command", nil, 5, time.Time{}, "", nil)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	is.True(comm != nil)

	comm.Status = "updated"

	err = comm.Save()
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	list, err := device.ListCommands(&dh.ListParams{
		Start: comm.Timestamp.Time.Add(-1 * time.Millisecond),
	})
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	is.True(len(list) > 0)
	is.Equal(list[len(list)-1].Status, "updated")

	err = device.Remove()
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}
}

func TestDeviceNotifications(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-notification", "", nil, 0, 0, false)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	notif, err := device.SendNotification("test notif", nil, time.Time{})
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	is.True(notif != nil)

	list, err := device.ListNotifications(nil)

	is.True(len(list) > 0)

	err = device.Remove()
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}
}
