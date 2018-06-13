// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package dh_wsclient_test

import (
	"encoding/json"
	"github.com/devicehive/devicehive-go"
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestWSClientSubscriptions(t *testing.T) {
	is := is.New(t)

	device := devicehive_go.Device{
		Id: "go-test-dev",
	}
	err := wsclient.PutDevice(device)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	defer func() {
		err := wsclient.DeleteDevice("go-test-dev")
		if err != nil {
			t.Fatal(err)
		}

		testResponse(t, nil)
	}()

	err = wsclient.SendDeviceCommand("go-test-dev", "go-test", nil, 120, time.Time{}, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	comm := &devicehive_go.Command{}
	testResponse(t, func(data []byte) {
		json.Unmarshal(data, comm)
	})

	err = wsclient.SubscribeCommands(&devicehive_go.SubscribeParams{
		DeviceId:  "go-test-dev",
		Timestamp: comm.Timestamp.Time.Add(-1 * time.Nanosecond),
	})
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)
	testResponse(t, func(data []byte) {
		res := make(map[string]json.RawMessage)
		insertedComm := &devicehive_go.Command{}
		json.Unmarshal(data, &res)
		json.Unmarshal(res["command"], insertedComm)
		is.Equal(insertedComm.Id, comm.Id)
	})

	err = wsclient.SendDeviceNotification("go-test-dev", "go-test", nil, time.Time{})
	if err != nil {
		t.Fatal(err)
	}

	notif := &devicehive_go.Notification{}
	testResponse(t, func(data []byte) {
		json.Unmarshal(data, notif)
	})

	err = wsclient.SubscribeNotifications(&devicehive_go.SubscribeParams{
		DeviceId:  "go-test-dev",
		Timestamp: notif.Timestamp.Time.Add(-1 * time.Millisecond),
	})
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)
	testResponse(t, func(data []byte) {
		res := make(map[string]json.RawMessage)
		insertedNotif := &devicehive_go.Notification{}
		json.Unmarshal(data, &res)
		json.Unmarshal(res["notification"], insertedNotif)
		is.Equal(insertedNotif.Id, notif.Id)
	})
}

func TestWSClientUnsubscriptions(t *testing.T) {
	err := wsclient.SubscribeCommands(nil)
	if err != nil {
		t.Fatal(err)
	}

	subsId := ""
	testResponse(t, func(data []byte) {
		res := make(map[string]interface{})
		json.Unmarshal(data, &res)
		subsId = res["subscriptionId"].(string)
	})

	wsclient.UnsubscribeCommands(subsId)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-wsclient.ErrorChan:
		t.Fatal(err)
	case <-time.After(TestTimeout):
	}

	err = wsclient.SubscribeNotifications(nil)
	if err != nil {
		t.Fatal(err)
	}

	subsId = ""
	testResponse(t, func(data []byte) {
		res := make(map[string]interface{})
		json.Unmarshal(data, &res)
		subsId = res["subscriptionId"].(string)
	})

	wsclient.UnsubscribeNotifications(subsId)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-wsclient.ErrorChan:
		t.Fatal(err)
	case <-time.After(TestTimeout):
	}
}
