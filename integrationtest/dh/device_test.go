package dh_test

import (
	"github.com/matryer/is"
	"testing"
	"time"
	"strconv"
	"github.com/devicehive/devicehive-go/dh"
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
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(len(list) > 0)
	is.Equal(list[len(list)-1].Status, "updated")

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

func TestDeviceSubscribeInsertCommands(t *testing.T) {
	waitTimeout := time.Duration(client.PollingWaitTimeoutSeconds+10) * time.Second

	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}
	defer func() {
		err := device.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	var lastCommand *dh.Command
	for i := int64(0); i < 5; i++ {
		lastCommand, err = device.SendCommand("go test command "+strconv.FormatInt(i, 10), nil, 120, time.Time{}, "", nil)
		if err != nil {
			t.Errorf("%s: %v", err.Name(), err)
			return
		}
	}

	commSubs, err := device.SubscribeInsertCommands(nil, lastCommand.Timestamp.Add(-3*time.Second))
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}
	defer func() {
		err := commSubs.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	for i := int64(0); i < 5; i++ {
		select {
		case comm, ok := <-commSubs.CommandsChan:
			is.True(ok)
			is.True(comm != nil)
			is.Equal(comm.Command, "go test command "+strconv.FormatInt(i, 10))
		case <-time.After(waitTimeout):
			t.Error("command insert event timeout")
		}
	}
}

func TestDeviceSubscribeUpdateCommands(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	commUpdChan, err := device.SubscribeUpdateCommands(nil, time.Time{})
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	go func() {
		select {
		case comm, ok := <-commUpdChan.CommandsChan:
			is.True(ok)
			is.True(comm != nil)
			is.Equal(comm.Status, "updated")
		case <-time.After(1 * time.Second):
			t.Error("command update event timeout")
		}

		err = device.Remove()
		if err != nil {
			t.Errorf("%s: %v", err.Name(), err)
			return
		}
	}()

	comm, err := device.SendCommand("go test command", nil, 5, time.Time{}, "", nil)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	comm.Status = "updated"

	err = comm.Save()
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	<-time.After(1 * time.Second)
}

func TestDeviceCommandSubscriptionRemove(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	commChan, err := device.SubscribeInsertCommands(nil, time.Time{})

	go func() {
		select {
		case comm, ok := <-commChan.CommandsChan:
			is.True(!ok)
			is.True(comm == nil)
		case <-time.After(300 * time.Millisecond):
			t.Error("command unsubscribe timeout")
		}

		err = device.Remove()
		if err != nil {
			t.Errorf("%s: %v", err.Name(), err)
			return
		}
	}()

	commChan.Remove()

	<-time.After(300 * time.Millisecond)
}

func TestDeviceSubscribeNotifications(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
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
			t.Errorf("%s: %v", err.Name(), err)
			return
		}
	}()

	_, err = device.SendNotification("go test notification", nil, time.Time{})
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	<-time.After(1 * time.Second)
}

func TestDeviceNotificationSubscriptionRemove(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	subs, err := device.SubscribeNotifications(nil, time.Time{})

	go func() {
		select {
		case comm, ok := <-subs.NotificationChan:
			is.True(!ok)
			is.True(comm == nil)
		case <-time.After(300 * time.Millisecond):
			t.Error("notification unsubscribe timeout")
		}

		err = device.Remove()
		if err != nil {
			t.Errorf("%s: %v", err.Name(), err)
			return
		}
	}()

	subs.Remove()

	<-time.After(300 * time.Millisecond)
}
