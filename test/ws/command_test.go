package dh_test

import (
	"testing"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/matryer/is"
	"github.com/devicehive/devicehive-go/dh"
	"time"
	"github.com/gorilla/websocket"
)

func TestCommandGet(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

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

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	listReqParams := &dh.ListParams{
		Start:        time.Now().Add(-1 * time.Hour),
		End:          time.Now(),
		Command: 	  "test command",
		Status:		  "created",
		SortField:    "timestamp",
		SortOrder:    "ASC",
		Take:         10,
		Skip:         5,
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

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	comm := &dh.Command{
		Timestamp: dh.ISO8601Time{time.Now()},
		Parameters: map[string]interface{} {
			"test": 1,
		},
		Lifetime: 120,
		Status: "created",
	}
	err = client.CommandInsert("device id", "name", comm)

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

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	comm := &dh.Command{
		Timestamp: dh.ISO8601Time{time.Now()},
		Parameters: map[string]interface{} {
			"test": 1,
		},
		Lifetime: 120,
		Status: "created",
	}
	err = client.CommandUpdate("device id",  111, comm)

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
		testTimeout = 1 * time.Second
	)

	wsTestSrv.SetRequestHandler(func(reqData map[string]interface{}, c *websocket.Conn) {
		subscribeResponse := stubs.ResponseStub.Respond(reqData)
		c.WriteJSON(subscribeResponse)
		<-time.After(commandInsertEventDelay)

		c.WriteJSON(stubs.ResponseStub.CommandInsertEvent(subscribeResponse["subscriptionId"], reqData["deviceId"]))
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
		ReturnUpdatedCommands: true,
		Limit: 100,
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
		is.True(comm.Id != 0)
		is.True(comm.Command != "")
		is.True(comm.Timestamp.Unix() > 0)
		is.Equal(comm.DeviceId, "device id")
		is.True(comm.Parameters != nil)
	case <-time.After(testTimeout):
		t.Error("comand insert event timeout")
	}
}
