package dh_test

import (
	"testing"
	"github.com/matryer/is"
	"github.com/gorilla/websocket"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/devicehive/devicehive-go/dh"
)

func TestNotificationGet(t *testing.T) {
	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, conn *websocket.Conn) map[string]interface{} {
		is.Equal(reqData["action"].(string), "notification/get")

		reqId := reqData["requestId"].(string)
		devId := reqData["deviceId"].(string)
		notifId := int64(reqData["notificationId"].(float64))
		return stubs.ResponseStub.NotificationGet(reqId, devId, notifId)
	})

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	notif, dhErr := client.NotificationGet("device id", 123456789)

	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		return
	}

	is.Equal(notif.Id, int64(123456789))
	is.True(notif.Name != "")
	is.True(notif.Timestamp.Unix() != 0)
	is.Equal(notif.DeviceId, "device id")
	is.True(notif.NetworkId != 0)
	is.True(notif.Parameters != nil)
}
