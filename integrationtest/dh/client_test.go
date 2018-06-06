// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package dh_test

import (
	dh "github.com/devicehive/devicehive-go"
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
	case <-time.After(2 * time.Second):
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

func TestCommandSubscriptionRemove(t *testing.T) {
	is := is.New(t)

	subs, err := client.SubscribeCommands(nil)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	err = subs.Remove()
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	select {
	case comm, ok := <-subs.CommandsChan:
		is.True(!ok)
		is.True(comm == nil)
	case <-time.After(300 * time.Millisecond):
		t.Fatalf("command unsubscribe timeout")
	}
}

func TestNotificationSubscriptionRemove(t *testing.T) {
	is := is.New(t)

	subs, err := client.SubscribeNotifications(nil)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

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
