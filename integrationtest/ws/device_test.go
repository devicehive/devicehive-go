package dh_test

import (
	"testing"
	"github.com/matryer/is"
)

func TestDevice(t *testing.T) {
	is := is.New(t)

	device, err := client.PutDevice("go-test-dev", "", nil, 0, 0, false)
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

	err = device.Remove()
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}
}
