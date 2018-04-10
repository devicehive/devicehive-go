package internal_test

import (
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/gorilla/websocket"
	"github.com/matryer/is"
	"testing"
	"time"
	"encoding/json"
)

const testTimeout = 300 * time.Millisecond

func TestRequestId(t *testing.T) {
	wsTestSrv := &stubs.WSTestServer{}

	addr := wsTestSrv.Start()
	defer wsTestSrv.Close()

	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		is.True(reqData["requestId"] != "")
		return reqData
	})

	wsTsp, err := transport.Create(addr)

	is.NoErr(err)

	wsTsp.Request(map[string]interface{}{}, testTimeout)
}

func TestTimeout(t *testing.T) {
	wsTestSrv := &stubs.WSTestServer{}

	addr := wsTestSrv.Start()
	defer wsTestSrv.Close()

	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		<-time.After(testTimeout + 1*time.Second)

		return map[string]interface{}{
			"result": "success",
		}
	})

	wsTsp, err := transport.Create(addr)

	is.NoErr(err)

	res, tspErr := wsTsp.Request(map[string]interface{}{}, testTimeout)

	is.True(res == nil)
	is.Equal(tspErr.Name(), transport.TimeoutErr)
}

func TestInvalidResponse(t *testing.T) {
	wsTestSrv := &stubs.WSTestServer{}

	addr := wsTestSrv.Start()
	defer wsTestSrv.Close()

	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		c.WriteMessage(websocket.TextMessage, []byte("invalid response"))
		return nil
	})

	wsTsp, err := transport.Create(addr)

	is.NoErr(err)

	res, tspErr := wsTsp.Request(map[string]interface{}{}, testTimeout)

	is.True(res == nil)
	is.Equal(tspErr.Name(), transport.TimeoutErr)
}

func TestConnectionClose(t *testing.T) {
	wsTestSrv := &stubs.WSTestServer{}

	addr := wsTestSrv.Start()
	defer wsTestSrv.Close()

	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		c.Close()
		return nil
	})

	wsTsp, err := transport.Create(addr)

	is.NoErr(err)

	res, tspErr := wsTsp.Request(map[string]interface{}{}, testTimeout)

	is.True(res == nil)
	is.Equal(tspErr.Name(), transport.ConnClosedErr)
}


func TestSubscribe(t *testing.T) {
	wsTestSrv := &stubs.WSTestServer{}

	addr := wsTestSrv.Start()
	defer wsTestSrv.Close()

	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		res := stubs.ResponseStub.Respond(reqData)

		c.WriteJSON(res)
		<- time.After(500 * time.Millisecond)
		c.WriteJSON(stubs.ResponseStub.NotificationInsertEvent(res["subscriptionId"], reqData["deviceId"]))

		return nil
	})

	wsTsp, err := transport.Create(addr)

	is.NoErr(err)

	res, err := wsTsp.Request(map[string]interface{}{
		"action": "notification/subscribe",
	}, 0)

	type subsId struct {
		Value int64 `json:"subscriptionId"`
	}
	sid := &subsId{}

	json.Unmarshal(res, sid)

	tspChan := wsTsp.Subscribe(sid.Value)

	select {
	case rawNotif, ok := <- tspChan:
		is.True(ok)
		is.True(rawNotif != nil)
	case <- time.After(1 * time.Second):
		t.Error("subscription event timeout")
	}
}