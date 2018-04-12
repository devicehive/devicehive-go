package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/gorilla/websocket"
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestNotificationGet(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	notif, err := client.NotificationGet("device id", 123456789)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.Equal(notif.Id, int64(123456789))
	is.True(notif.Notification != "")
	is.True(notif.Timestamp.Unix() > 0)
	is.Equal(notif.DeviceId, "device id")
	is.True(notif.NetworkId != 0)
	is.True(notif.Parameters != nil)
}

func TestNotificationList(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	listReqParams := &dh.ListParams{
		Start:        time.Now().Add(-1 * time.Hour),
		End:          time.Now(),
		Notification: "test notif",
		SortField:    "timestamp",
		SortOrder:    "ASC",
		Take:         10,
		Skip:         5,
	}
	list, err := client.NotificationList("device id", listReqParams)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(len(list) != 0)
}

func TestNotificationInsert(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	devId := "device id"
	name := "test notif"
	ts := time.Now()
	params := map[string]interface{}{
		"testParam": 1,
	}
	notifId, err := client.NotificationInsert(devId, name, ts, params)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(notifId != 0)
}

func TestNotificationSubscribe(t *testing.T) {
	wsTestSrv, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	const (
		notifInsertEventDelay = 200 * time.Millisecond
		testTimeout = 1 * time.Second
	)

	wsTestSrv.SetRequestHandler(func(reqData map[string]interface{}, c *websocket.Conn) {
		subscribeResponse := stubs.ResponseStub.Respond(reqData)
		c.WriteJSON(subscribeResponse)
		<-time.After(notifInsertEventDelay)

		c.WriteJSON(stubs.ResponseStub.NotificationInsertEvent(subscribeResponse["subscriptionId"], reqData["deviceId"]))
	})

	is := is.New(t)

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	subsParams := &dh.SubscribeParams{
		Timestamp:     time.Now(),
		DeviceId:      "device id",
		NetworkIds:    []string{"net1", "net2"},
		DeviceTypeIds: []string{"dt1", "dt2"},
		Names:         []string{"n1", "n2"},
	}
	notifChan, err := client.NotificationSubscribe(subsParams)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(notifChan != nil)

	select {
	case notif, ok := <-notifChan:
		is.True(ok)
		is.True(notif.Id != 0)
		is.True(notif.Notification != "")
		is.True(notif.Timestamp.Unix() > 0)
		is.Equal(notif.DeviceId, "device id")
		is.True(notif.Parameters != nil)
	case <-time.After(testTimeout):
		t.Error("notification insert event timeout")
	}
}

func TestNotificationUnsubscribe(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	notifChan, err := client.NotificationSubscribe(nil)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	err = client.NotificationUnsubscribe(notifChan)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}
}
