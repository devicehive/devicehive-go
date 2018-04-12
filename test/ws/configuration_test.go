package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/matryer/is"
	"testing"
)

func TestConfigurationGet(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	conf, dhErr := client.ConfigurationGet("some_config")

	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		t.Fail()
	}

	is.True(conf != nil)
	is.Equal(conf.Name, "some_config")
	is.True(conf.Value != "")
}

func TestConfigurationPut(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	conf, dhErr := client.ConfigurationPut("some_config", "some test value")

	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		t.Fail()
	}

	is.True(conf != nil)
	is.Equal(conf.Name, "some_config")
	is.True(conf.Value == "some test value")
}

func TestConfigurationDelete(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	dhErr := client.ConfigurationDelete("some_config")

	if dhErr != nil {
		t.Error(dhErr)
		t.Fail()
	}
}
