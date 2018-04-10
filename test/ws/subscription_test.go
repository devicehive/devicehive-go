package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/matryer/is"
	"testing"
	"github.com/devicehive/devicehive-go/test/stubs"
)

func TestSubscriptionList(t *testing.T) {
	wsTestSrv := &stubs.WSTestServer{}

	addr := wsTestSrv.Start()
	defer wsTestSrv.Close()

	is := is.New(t)

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	subscriptions, dhErr := client.SubscriptionList(dh.NotificationType)

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
	is.True(subscriptions[0].Timestamp.Unix() > 0)

	is.True(subscriptions[1].NetworkIds == nil)
	is.True(subscriptions[1].DeviceTypeIds == nil)
	is.True(subscriptions[1].Names == nil)
}
