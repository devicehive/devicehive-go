package dh_test

import (
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/matryer/is"
	"testing"
)

func TestGetProperty(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client := connect(addr)

	prop, dhErr := client.GetProperty("some.test.prop")

	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		t.Fail()
	}

	is.True(prop != nil)
	is.Equal(prop.Name, "some.test.prop")
	is.True(prop.Value != "")
}

func TestSetProperty(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client := connect(addr)

	entityVersion, dhErr := client.SetProperty("some.test.prop", "some test value")

	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		t.Fail()
	}

	is.True(entityVersion != -1)
}

func TestDeleteProperty(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	client := connect(addr)

	dhErr := client.DeleteProperty("some.test.prop")

	if dhErr != nil {
		t.Error(dhErr)
		t.Fail()
	}
}
