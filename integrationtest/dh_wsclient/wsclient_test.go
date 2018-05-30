package dh_wsclient_test

import (
	"testing"
	"time"
	"github.com/devicehive/devicehive-go"
	"encoding/json"
	"github.com/matryer/is"
)

func TestWSClientSubscriptions(t *testing.T) {
	is := is.New(t)

	err := wsclient.PutDevice("go-test-dev", "", nil, 0, 0, false)
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
		Timestamp: notif.Timestamp.Time.Add(-500 * time.Millisecond),
	})
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, func(data []byte) {
		res := make(map[string]json.RawMessage)
		insertedNotif := &devicehive_go.Notification{}
		json.Unmarshal(data, &res)
		json.Unmarshal(res["notification"], insertedNotif)
		is.Equal(insertedNotif.Id, notif.Id)
	})
}
