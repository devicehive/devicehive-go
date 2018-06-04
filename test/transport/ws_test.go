package transport

import (
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/devicehive/devicehive-go/transport"
	"github.com/gorilla/websocket"
	"github.com/matryer/is"
	"testing"
	"time"
)

const testWSTimeout = 300 * time.Millisecond

func TestWSRequestId(t *testing.T) {
	wsTestSrv, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	wsTestSrv.SetRequestHandler(func(reqData map[string]interface{}, c *websocket.Conn) {
		is.True(reqData["requestId"] != "")
		c.WriteJSON(reqData)
	})

	wsTsp, err := transport.Create(addr)

	is.NoErr(err)

	wsTsp.Request("", nil, testWSTimeout)
}

func TestWSTimeout(t *testing.T) {
	wsTestSrv, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	wsTestSrv.SetRequestHandler(func(reqData map[string]interface{}, c *websocket.Conn) {
		<-time.After(testWSTimeout + 1*time.Second)

		c.WriteJSON(map[string]interface{}{
			"result": "success",
		})
	})

	wsTsp, err := transport.Create(addr)

	is.NoErr(err)

	res, tspErr := wsTsp.Request("", nil, testWSTimeout)

	is.True(res == nil)
	is.Equal(tspErr.Name(), transport.TimeoutErr)
}

func TestWSInvalidResponse(t *testing.T) {
	wsTestSrv, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	wsTestSrv.SetRequestHandler(func(reqData map[string]interface{}, c *websocket.Conn) {
		c.WriteMessage(websocket.TextMessage, []byte("invalid response"))
	})

	wsTsp, err := transport.Create(addr)

	is.NoErr(err)

	res, tspErr := wsTsp.Request("", nil, testWSTimeout)

	is.True(res == nil)
	is.Equal(tspErr.Name(), transport.TimeoutErr)
}

func TestWSConnectionClose(t *testing.T) {
	wsTestSrv, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	wsTestSrv.SetRequestHandler(func(reqData map[string]interface{}, c *websocket.Conn) {
		c.Close()
	})

	wsTsp, err := transport.Create(addr)

	is.NoErr(err)

	res, tspErr := wsTsp.Request("", nil, testWSTimeout)

	is.True(res == nil)
	is.Equal(tspErr.Name(), transport.ConnClosedErr)
}

func TestWSSubscribe(t *testing.T) {
	wsTestSrv, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	wsTestSrv.SetRequestHandler(func(reqData map[string]interface{}, c *websocket.Conn) {
		res := stubs.ResponseStub.Respond(reqData)

		c.WriteJSON(stubs.ResponseStub.NotificationInsertEvent(res["subscriptionId"], reqData["deviceId"]))
		c.WriteJSON(res)
		c.WriteJSON(stubs.ResponseStub.NotificationInsertEvent(res["subscriptionId"], reqData["deviceId"]))
	})

	wsTsp, err := transport.Create(addr)

	is.NoErr(err)

	tspChan, _, tspErr := wsTsp.Subscribe("notification/subscribe", nil)
	if tspErr != nil {
		t.Fatalf("%s: %v", tspErr.Name(), tspErr)
	}

	select {
	case rawNotif, ok := <-tspChan:
		is.True(ok)
		is.True(rawNotif != nil)
	case <-time.After(1 * time.Second):
		t.Fatal("subscription event timeout")
	}

	select {
	case rawNotif, ok := <-tspChan:
		is.True(ok)
		is.True(rawNotif != nil)
	case <-time.After(1 * time.Second):
		t.Fatal("subscription event timeout")
	}
}
