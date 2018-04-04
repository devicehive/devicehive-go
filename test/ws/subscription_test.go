package dh_ws_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/gorilla/websocket"
	"github.com/matryer/is"
	"testing"
)

func TestSubscriptionList(t *testing.T) {
	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, conn *websocket.Conn) map[string]interface{} {
		return stubs.ResponseStub.SubscriptionList(reqData["requestId"].(string), reqData["type"].(string))
	})

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	subscriptions, dhErr := client.SubscriptionList(dh.Notification)

	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
	}

	is.True(subscriptions != nil)
	is.True(len(subscriptions) != 0)
	is.True(subscriptions[0].Id != 0)
	is.True(subscriptions[0].Type != "")
	is.True(subscriptions[0].DeviceId != "")
	is.True(subscriptions[0].NetworkIds != nil)
	is.True(subscriptions[0].DeviceTypeIds != nil)
	is.True(subscriptions[0].Names != nil)
	is.True(subscriptions[0].Timestamp.Unix() != 0)

	is.True(subscriptions[1].NetworkIds == nil)
	is.True(subscriptions[1].DeviceTypeIds == nil)
	is.True(subscriptions[1].Names == nil)
}
