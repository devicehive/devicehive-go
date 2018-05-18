package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/matryer/is"
	"testing"
	"time"
)

const commandsToSend = 1
const notificationsToSend = 1

const subscriptionTimestampOffset = -1 * time.Second

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
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev-comm-insert", "", nil, 0, 0, false)
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

	commSubs, err := device.SubscribeInsertCommands([]string{"go test command"}, firstValidCommand.Timestamp.Add(subscriptionTimestampOffset))
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

func TestDeviceSubscribeUpdateCommands(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev-comm-update", "", nil, 0, 0, false)
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
		comm, err := device.SendCommand("go test command", nil, 5, time.Time{}, "", nil)
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}

		comm.Status = "updated"

		err = comm.Save()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}

		if firstValidCommand == nil {
			firstValidCommand = comm
		}
	}

	commUpdSubs, err := device.SubscribeUpdateCommands(nil, firstValidCommand.Timestamp.Add(subscriptionTimestampOffset))
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}
	defer func() {
		err := commUpdSubs.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	for i := 0; i < commandsToSend; i++ {
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
	defer func() {
		err := device.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	commSubs, err := device.SubscribeInsertCommands(nil, time.Time{})

	err = commSubs.Remove()
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	select {
	case comm, ok := <-commSubs.CommandsChan:
		is.True(!ok)
		is.True(comm == nil)
	case <-time.After(300 * time.Millisecond):
		t.Fatalf("command unsubscribe timeout")
	}
}

func TestDeviceSubscribeNotifications(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev-notif-insert", "", nil, 0, 0, false)
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

	notifSubs, err := device.SubscribeNotifications(nil, firstValidNotif.Timestamp.Add(subscriptionTimestampOffset))
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

func TestDeviceNotificationSubscriptionRemove(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	subs, err := device.SubscribeNotifications(nil, time.Time{})
	defer func() {
		err = device.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	err = subs.Remove()
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	select {
	case comm, ok := <-subs.NotificationChan:
		is.True(!ok)
		is.True(comm == nil)
	case <-time.After(300 * time.Millisecond):
		t.Fatalf("notification unsubscribe timeout")
	}
}
