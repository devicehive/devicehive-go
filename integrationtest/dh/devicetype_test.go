package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestDeviceType(t *testing.T) {
	is := is.New(t)

	deviceType, err := client.CreateDeviceType("go-test-dev", "go sdk test device type")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = deviceType.Remove()
		if err != nil {
			t.Fatal(err)
		}
	}()

	is.True(deviceType != nil)
	is.True(deviceType.Id != 0)

	deviceType.Description = "updated go sdk test network"
	err = deviceType.Save()
	if err != nil {
		t.Fatal(err)
	}

	sameDeviceType, err := client.GetDeviceType(deviceType.Id)
	if err != nil {
		t.Fatal(err)
	}

	is.True(sameDeviceType != nil)
	is.Equal(sameDeviceType.Name, "go-test-dev")

	list, err := client.ListDeviceTypes(&dh.ListParams{
		NamePattern: "go-%-dev",
	})
	if err != nil {
		t.Fatal(err)
	}

	is.Equal(len(list), 1)
	is.Equal(list[0].Name, "go-test-dev")
}

func TestDeviceTypeSubscribeInsertCommands(t *testing.T) {
	is := is.New(t)

	devType, err := client.CreateDeviceType("go-test", "")
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}
	defer func() {
		err := devType.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	device, err := client.PutDevice("go-test-dev-comm-insert", "", nil, 0, devType.Id, false)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}
	defer func() {
		err := device.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	var firstValidCommand *dh.Command
	for i := 0; i < commandsToSend; i++ {
		comm, err := device.SendCommand("go test command", nil, 120, time.Time{}, "", nil)
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}

		if firstValidCommand == nil {
			firstValidCommand = comm
		}

		_, err = device.SendCommand("go test command to omit", nil, 120, time.Time{}, "", nil)
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}

	commSubs, err := devType.SubscribeInsertCommands(nil, firstValidCommand.Timestamp.Add(subscriptionTimestampOffset))
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}
	defer func() {
		err := commSubs.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	for i := 0; i < commandsToSend; i++ {
		select {
		case comm, ok := <-commSubs.CommandsChan:
			is.True(ok)
			is.True(comm != nil)
			is.Equal(comm.Command, "go test command")
		case <-time.After(waitTimeout):
			t.Error("command insert event timeout")
		}
	}
}

func TestDeviceTypeSubscribeNotifications(t *testing.T) {
	is := is.New(t)

	devType, err := client.CreateDeviceType("go-test", "")
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}
	defer func() {
		err := devType.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	device, err := client.PutDevice("go-test-dev-notif-insert", "", nil, 0, devType.Id, false)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}
	defer func() {
		err := device.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	var firstValidNotif *dh.Notification
	for i := 0; i < notificationsToSend; i++ {
		notif, err := device.SendNotification("go test notification", nil, time.Time{})
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}

		if firstValidNotif == nil {
			firstValidNotif = notif
		}
	}

	notifSubs, err := devType.SubscribeNotifications(nil, firstValidNotif.Timestamp.Add(subscriptionTimestampOffset))
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}
	defer func() {
		err := notifSubs.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	for i := 0; i < notificationsToSend; i++ {
		select {
		case notif, ok := <-notifSubs.NotificationChan:
			is.True(ok)
			is.True(notif != nil)
			is.Equal(notif.Notification, "go test notification")
		case <-time.After(waitTimeout):
			t.Error("notification insert event timeout")
		}
	}
}
