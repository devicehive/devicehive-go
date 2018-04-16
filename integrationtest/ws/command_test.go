package dh_test

import (
	"testing"
	"github.com/matryer/is"
	"github.com/devicehive/devicehive-go/dh"
	"time"
)

func TestCommand(t *testing.T) {
	err := auth()

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is := is.New(t)

	devId := "4NemW3PE9BHRSqb0DVVgsphZh7SCZzgm3Lxg"
	commData := &dh.Command{
		Timestamp: dh.ISO8601Time{time.Now()},
		Parameters: map[string]interface{} {
			"test": 1,
		},
		Lifetime: 5,
		Status: "created",
	}
	err = client.CommandInsert(devId, "name", commData)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(commData.Id != 0)

	commUpdate := &dh.Command{
		Status: "updated",
	}
	err = client.CommandUpdate(devId, commData.Id, commUpdate)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	comm, err := client.CommandGet(devId, commData.Id)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.Equal(comm.Id, commData.Id)
	is.Equal(comm.Status, commUpdate.Status)

	list, err := client.CommandList(devId, nil)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(len(list) != 0)
}

func TestCommandInsertSubscribe(t *testing.T) {
	err := auth()

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is := is.New(t)

	devId := "4NemW3PE9BHRSqb0DVVgsphZh7SCZzgm3Lxg"
	name := "test command"

	commChan, err := client.CommandSubscribe(nil)

	go func() {
		select {
		case comm := <-commChan:
			is.Equal(comm.Command, name)
		case <-time.After(1 * time.Second):
			t.Error("command insert event timeout")
		}
	}()

	comm :=&dh.Command{
		Lifetime: 5,
	}
	err = client.CommandInsert(devId, name, comm)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	<-time.After(500 * time.Millisecond)
}

func TestCommandUpdateSubscribe(t *testing.T) {
	err := auth()

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is := is.New(t)

	devId := "4NemW3PE9BHRSqb0DVVgsphZh7SCZzgm3Lxg"
	name := "test command"

	params := &dh.SubscribeParams{
		ReturnUpdatedCommands: true,
	}
	commChan, err := client.CommandSubscribe(params)

	go func() {
		<-commChan

		select {
		case comm := <-commChan:
			is.Equal(comm.Status, "updated")
		case <-time.After(1 * time.Second):
			t.Error("command insert event timeout")
		}
	}()

	comm :=&dh.Command{
		Lifetime: 5,
	}
	err = client.CommandInsert(devId, name, comm)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	upd := &dh.Command{
		Status: "updated",
	}
	err = client.CommandUpdate(devId, comm.Id, upd)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	<-time.After(500 * time.Millisecond)
}
