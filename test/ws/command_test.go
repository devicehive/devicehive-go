package dh_test

import (
	"testing"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/matryer/is"
	"github.com/devicehive/devicehive-go/dh"
	"time"
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
