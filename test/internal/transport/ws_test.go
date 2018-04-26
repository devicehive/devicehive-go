package transport_test

import (
	"encoding/json"
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/gorilla/websocket"
	"github.com/matryer/is"
	"strconv"
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

	wsTsp.Request("", map[string]interface{}{}, testWSTimeout)
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

	res, tspErr := wsTsp.Request("", map[string]interface{}{}, testWSTimeout)

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

	res, tspErr := wsTsp.Request("", map[string]interface{}{}, testWSTimeout)

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

	res, tspErr := wsTsp.Request("", map[string]interface{}{}, testWSTimeout)

	is.True(res == nil)
	is.Equal(tspErr.Name(), transport.ConnClosedErr)
}

func TestWSSubscribe(t *testing.T) {
	wsTestSrv, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	wsTestSrv.SetRequestHandler(func(reqData map[string]interface{}, c *websocket.Conn) {
		res := stubs.ResponseStub.Respond(reqData)

		c.WriteJSON(res)
		<-time.After(500 * time.Millisecond)
		c.WriteJSON(stubs.ResponseStub.NotificationInsertEvent(res["subscriptionId"], reqData["deviceId"]))
	})

	wsTsp, err := transport.Create(addr)

	is.NoErr(err)

	res, tspErr := wsTsp.Request("notification/subscribe", map[string]interface{}{}, 0)

	if tspErr != nil {
		t.Errorf("%s: %v", tspErr.Name(), tspErr)
		return
	}

	type subsId struct {
		Value int64 `json:"subscriptionId"`
	}
	sid := &subsId{}

	json.Unmarshal(res, sid)

	tspChan := wsTsp.Subscribe(strconv.FormatInt(sid.Value, 10))

	select {
	case rawNotif, ok := <-tspChan:
		is.True(ok)
		is.True(rawNotif != nil)
	case <-time.After(1 * time.Second):
		t.Error("subscription event timeout")
	}
}

func TestWSUnsubscribe(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	wsTsp, err := transport.Create(addr)

	is.NoErr(err)

	res, tspErr := wsTsp.Request("notification/subscribe", map[string]interface{}{}, 0)

	if tspErr != nil {
		t.Errorf("%s: %v", tspErr.Name(), tspErr)
		return
	}

	type subsId struct {
		Value int64 `json:"subscriptionId"`
	}
	sid := &subsId{}

	json.Unmarshal(res, sid)

	sidStr := strconv.FormatInt(sid.Value, 10)

	wsTsp.Subscribe(sidStr)

	wsTsp.Unsubscribe(sidStr)
}