package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/gorilla/websocket"
	"github.com/matryer/is"
	"testing"
)

func TestConfigurationGet(t *testing.T) {
	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, conn *websocket.Conn) map[string]interface{} {
		is.Equal(reqData["action"].(string), "configuration/get")
		return stubs.ResponseStub.ConfigurationGet(reqData["requestId"].(string), reqData["name"].(string))
	})

	client, err := dh.Connect(wsServerAddr)

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
	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, conn *websocket.Conn) map[string]interface{} {
		is.Equal(reqData["action"].(string), "configuration/put")
		return stubs.ResponseStub.ConfigurationPut(reqData["requestId"].(string), reqData["name"].(string), reqData["value"].(string))
	})

	client, err := dh.Connect(wsServerAddr)

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
	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, conn *websocket.Conn) map[string]interface{} {
		is.Equal(reqData["action"].(string), "configuration/delete")
		return stubs.ResponseStub.EmptySuccessResponse("configuration/delete", reqData["requestId"].(string))
	})

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	dhErr := client.ConfigurationDelete("some_config")

	if dhErr != nil {
		t.Error(dhErr)
		t.Fail()
	}
}
