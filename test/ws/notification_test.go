package dh_test

import (
	"testing"
	"github.com/matryer/is"
	"github.com/devicehive/devicehive-go/dh"
	"time"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/gorilla/websocket"
)

func TestNotificationGet(t *testing.T) {
	wsTestSrv := &stubs.WSTestServer{}

	addr := wsTestSrv.Start()
	defer wsTestSrv.Close()

	is := is.New(t)

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	notif, dhErr := client.NotificationGet("device id", 123456789)
	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
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
	wsTestSrv := &stubs.WSTestServer{}

	addr := wsTestSrv.Start()
	defer wsTestSrv.Close()

	is := is.New(t)

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	listReqParams := &dh.ListParams{
		Start: time.Now().Add(-1 * time.Hour),
		End: time.Now(),
		Notification: "test notif",
		SortField: "timestamp",
		SortOrder: "ASC",
		Take: 10,
		Skip: 5,
	}
	list, dhErr := client.NotificationList("device id", listReqParams)

	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		return
	}

	is.True(len(list) != 0)
}

func TestNotificationInsert(t *testing.T) {
	wsTestSrv := &stubs.WSTestServer{}

	addr := wsTestSrv.Start()
	defer wsTestSrv.Close()

	is := is.New(t)

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	devId := "device id"
	name := "test notif"
	ts := time.Now()
	params := map[string]interface{} {
		"testParam": 1,
	}
	notifId, dhErr := client.NotificationInsert(devId, name, ts, params)
	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		return
	}

	is.True(notifId != 0)
}

func TestNotificationSubscribe(t *testing.T) {
	wsTestSrv := &stubs.WSTestServer{}

	addr := wsTestSrv.Start()
	defer wsTestSrv.Close()

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, conn *websocket.Conn) map[string]interface{} {
		res := stubs.ResponseStub.Respond(reqData)
		conn.WriteJSON(res)
		<- time.After(200 * time.Millisecond)

		return stubs.ResponseStub.NotificationInsertEvent(res["subscriptionId"], reqData["deviceId"])
	})

	is := is.New(t)

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	subsParams := &dh.SubscribeParams{
		Timestamp: time.Now(),
		DeviceId: "device id",
		NetworkIds: []string{ "net1", "net2" },
		DeviceTypeIds: []string{ "dt1", "dt2" },
		Names: []string{ "n1", "n2" },
	}
	notifChan, dhErr := client.NotificationSubscribe(subsParams)
	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		return
	}

	is.True(notifChan != nil)

	select {
	case notif, ok := <- notifChan:
		is.True(ok)
		is.True(notif.Id != 0)
		is.True(notif.Notification != "")
		is.True(notif.Timestamp.Unix() > 0)
		is.Equal(notif.DeviceId, "device id")
		is.True(notif.Parameters != nil)
	case <- time.After(1 * time.Second):
		t.Error("notification insert event timeout")
	}
}
