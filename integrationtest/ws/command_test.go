package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestCommand(t *testing.T) {
	is := is.New(t)

	commData := &dh.Command{
		Timestamp: dh.ISO8601Time{time.Now()},
		Parameters: map[string]interface{}{
			"test": 1,
		},
		Lifetime: 5,
		Status:   "created",
	}
	err := client.CommandInsert(testDeviceId, "name", commData)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(commData.Id != 0)

	commUpdate := &dh.Command{
		Status: "updated",
	}
	err = client.CommandUpdate(testDeviceId, commData.Id, commUpdate)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	comm, err := client.CommandGet(testDeviceId, commData.Id)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.Equal(comm.Id, commData.Id)
	is.Equal(comm.Status, commUpdate.Status)

	list, err := client.CommandList(testDeviceId, nil)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(len(list) != 0)
}

func TestCommandInsertSubscribe(t *testing.T) {
	is := is.New(t)

	name := "test command insert"

	commChan, err := client.CommandSubscribe(nil)

	go func() {
		select {
		case comm := <-commChan:
			is.Equal(comm.Command, name)
		case <-time.After(1 * time.Second):
			t.Error("command insert event timeout")
		}
	}()

	comm := &dh.Command{
		Lifetime: 5,
	}
	err = client.CommandInsert(testDeviceId, name, comm)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	<-time.After(500 * time.Millisecond)
}

// @TODO This test is impacted by some other test, so it fails when launched with all test suit, but passes when launched alone
func TestCommandUpdateSubscribe(t *testing.T) {
	is := is.New(t)

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
			t.Error("command update event timeout")
		}
	}()

	name := "test command update"
	comm := &dh.Command{
		Lifetime: 5,
	}
	err = client.CommandInsert(testDeviceId, name, comm)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	upd := &dh.Command{
		Status: "updated",
	}
	err = client.CommandUpdate(testDeviceId, comm.Id, upd)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	<-time.After(500 * time.Millisecond)
}

func TestCommandUnsubscribe(t *testing.T) {
	commChan, err := client.CommandSubscribe(nil)

	go func() {
		select {
		case comm, ok := <-commChan:
			if comm != nil || ok {
				t.Error("client hasn't been unsubscribed")
			}
		case <-time.After(1 * time.Second):
			t.Error("timeout")
		}
	}()

	err = client.CommandUnsubscribe(commChan)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	name := "test command"
	client.CommandInsert(testDeviceId, name, nil)
}
