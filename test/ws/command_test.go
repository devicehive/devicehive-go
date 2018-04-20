package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/gorilla/websocket"
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestCommandGet(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client := connect(addr)

	comm, err := client.CommandGet("device id", 1111111)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.Equal(comm.Id, int64(1111111))
	is.Equal(comm.DeviceId, "device id")
}

func TestCommandList(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client := connect(addr)

	listReqParams := &dh.ListParams{
		Start:     time.Now().Add(-1 * time.Hour),
		End:       time.Now(),
		Command:   "test command",
		Status:    "created",
		SortField: "timestamp",
		SortOrder: "ASC",
		Take:      10,
		Skip:      5,
	}
	list, err := client.CommandList("device id", listReqParams)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(len(list) != 0)
}

func TestCommandInsert(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client := connect(addr)

	comm := &dh.Command{
		Timestamp: dh.ISO8601Time{time.Now()},
		Parameters: map[string]interface{}{
			"test": 1,
		},
		Lifetime: 120,
		Status:   "created",
	}
	err := client.CommandInsert("device id", "name", comm)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.Equal(comm.DeviceId, "device id")
	is.Equal(comm.Command, "name")
	is.True(comm.Id != 0)
	is.True(comm.LastUpdated.Unix() > 0)
	is.True(comm.UserId != 0)
}

func TestCommandUpdate(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	client := connect(addr)

	comm := &dh.Command{
		Timestamp: dh.ISO8601Time{time.Now()},
		Parameters: map[string]interface{}{
			"test": 1,
		},
		Lifetime: 120,
		Status:   "created",
	}
	err := client.CommandUpdate("device id", 111, comm)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}
}

func TestCommandSubscribe(t *testing.T) {
	wsTestSrv, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	const (
		commandInsertEventDelay = 200 * time.Millisecond
		testTimeout             = 1 * time.Second
	)

	wsTestSrv.SetRequestHandler(func(reqData map[string]interface{}, c *websocket.Conn) {
		subscribeResponse := stubs.ResponseStub.Respond(reqData)
		c.WriteJSON(subscribeResponse)
		<-time.After(commandInsertEventDelay)

		c.WriteJSON(stubs.ResponseStub.CommandInsertEvent(subscribeResponse["subscriptionId"], reqData["deviceId"]))
	})

	is := is.New(t)

	client := connect(addr)

	subsParams := &dh.SubscribeParams{
		Timestamp: time.Now(),
		DeviceId:  "device id",
		Limit:     100,
	}
	commChan, err := client.CommandSubscribe(subsParams)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(commChan != nil)

	select {
	case comm, ok := <-commChan:
		is.True(ok)
		is.True(comm != nil)
	case <-time.After(testTimeout):
		t.Error("comand insert event timeout")
	}
}

func TestCommandUnsubscribe(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	client := connect(addr)

	is := is.New(t)

	commChan, err := client.CommandSubscribe(nil)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	err = client.CommandUnsubscribe(commChan)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	n, ok := <-commChan

	is.True(n == nil)
	is.Equal(ok, false)
}

func TestCommandSave(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	client := connect(addr)

	device, err := client.GetDevice("device-id")
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	list, err := device.ListCommands(nil)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	comm := list[0]

	comm.Status = "updated"

	err = comm.Save()
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}
}
