package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestSubscriptions(t *testing.T) {
	is := is.New(t)
	device, err := client.PutDevice("go-test-subs", "", nil, 0, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := device.Remove()
		if err != nil {
			t.Fatalf("%s: %v", err.Name(), err)
		}
	}()

	notif, err := device.SendNotification("test", nil, time.Time{})
	if err != nil {
		t.Fatal(err)
	}

	notifSubs, err := client.SubscribeNotifications(&dh.SubscribeParams{
		DeviceId:  "go-test-subs",
		Timestamp: notif.Timestamp.Time.Add(-1 * time.Nanosecond),
	})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case n := <-notifSubs.NotificationChan:
		is.Equal(n.Notification, "test")
	case <-time.After(1 * time.Second):
		t.Fatal("Subscription event timeout")
	}

	comm, err := device.SendCommand("test", nil, 0, time.Time{}, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	commSubs, err := client.SubscribeCommands(&dh.SubscribeParams{
		DeviceId:  "go-test-subs",
		Timestamp: comm.Timestamp.Time.Add(-1 * time.Nanosecond),
	})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case c := <-commSubs.CommandsChan:
		is.Equal(c.Command, "test")
	case <-time.After(1 * time.Second):
		t.Fatal("Subscription event timeout")
	}
}
