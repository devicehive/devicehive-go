package dh_test

import (
	"testing"
	"github.com/matryer/is"
	"github.com/devicehive/devicehive-go/dh"
)

func TestDevice(t *testing.T) {
	is := is.New(t)

	device := &dh.Device{
		Name: "go test",
	}

	err := client.PutDevice("go-test-dev", device)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	device.Name = "go test updated"

	err = client.PutDevice(device.Id, device)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	device, err = client.GetDevice(device.Id)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(device != nil)
	is.Equal(device.Name, "go test updated")

	err = client.RemoveDevice("go-test-dev")
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}
}
