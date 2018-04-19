package dh_test

import (
	"testing"
	"github.com/matryer/is"
)

func TestDevice(t *testing.T) {
	is := is.New(t)

	err := client.PutDevice("go-test-dev", nil)
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	device, err := client.GetDevice("go-test-dev")
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is.True(device != nil)

	err = client.RemoveDevice("go-test-dev")
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}
}
