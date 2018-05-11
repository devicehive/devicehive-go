package dh_test

import (
	"github.com/matryer/is"
	"testing"
	"time"
	"github.com/devicehive/devicehive-go/dh"
)

func TestDevice(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

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

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
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

	list, err := device.ListCommands(nil)
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

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
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

func TestDeviceSubscribeInsertCommands(t *testing.T) {
	const commandsCount = 5
	waitTimeout := time.Duration(client.PollingWaitTimeoutSeconds+10) * time.Second

	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}
	defer func() {
		err := device.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	var comm *dh.Command
	for i := int64(0); i < commandsCount; i++ {
		comm, err = device.SendCommand("go test command", nil, 120, time.Time{}, "", nil)
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}

	commSubs, err := device.SubscribeInsertCommands(nil, comm.Timestamp.Add(-3*time.Second))
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}
	defer func() {
		err := commSubs.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	for i := int64(0); i < commandsCount; i++ {
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

func TestDeviceSubscribeUpdateCommands(t *testing.T) {
	const commandsCount = 5
	waitTimeout := time.Duration(client.PollingWaitTimeoutSeconds+10) * time.Second

	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}
	defer func() {
		err := device.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	var comm *dh.Command
	for i := 0; i < commandsCount; i++ {
		comm, err = device.SendCommand("go test command", nil, 5, time.Time{}, "", nil)
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}

		comm.Status = "updated"

		err = comm.Save()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}

	commUpdSubs, err := device.SubscribeUpdateCommands(nil, comm.Timestamp.Add(-3 * time.Second))
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}
	defer func() {
		err := commUpdSubs.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	for i := 0; i < commandsCount; i++ {
		select {
		case comm, ok := <-commUpdSubs.CommandsChan:
			is.True(ok)
			is.True(comm != nil)
			is.Equal(comm.Status, "updated")
		case <-time.After(waitTimeout):
			t.Fatal("command update event timeout")
		}
	}
}

func TestDeviceCommandSubscriptionRemove(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	commChan, err := device.SubscribeInsertCommands(nil, time.Time{})

	go func() {
		select {
		case comm, ok := <-commChan.CommandsChan:
			is.True(!ok)
			is.True(comm == nil)
		case <-time.After(300 * time.Millisecond):
			t.Fatalf("command unsubscribe timeout")
		}

		err = device.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	commChan.Remove()

	<-time.After(300 * time.Millisecond)
}

func TestDeviceSubscribeNotifications(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	notifChan, err := device.SubscribeNotifications(nil, time.Time{})

	go func() {
		select {
		case notif, ok := <-notifChan.NotificationChan:
			is.True(ok)
			is.True(notif != nil)
			is.Equal(notif.Notification, "go test notification")
		case <-time.After(1 * time.Second):
			t.Error("notification insert event timeout")
		}

		err = device.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	_, err = device.SendNotification("go test notification", nil, time.Time{})
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	<-time.After(1 * time.Second)
}

func TestDeviceNotificationSubscriptionRemove(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	subs, err := device.SubscribeNotifications(nil, time.Time{})

	go func() {
		select {
		case comm, ok := <-subs.NotificationChan:
			is.True(!ok)
			is.True(comm == nil)
		case <-time.After(300 * time.Millisecond):
			t.Fatalf("notification unsubscribe timeout")
		}

		err = device.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	subs.Remove()

	<-time.After(300 * time.Millisecond)
}
