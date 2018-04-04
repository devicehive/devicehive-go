package dh_ws_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/gorilla/websocket"
	"github.com/matryer/is"
	"testing"
)

func TestServerInfo(t *testing.T) {
	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		is.Equal(reqData["action"].(string), "server/info")
		return stubs.ResponseStub.ServerInfo(reqData["requestId"].(string))
	})

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	res, dhErr := client.ServerInfo()

	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
	}

	is.True(res != nil)
	is.True(res.APIVersion != "")
	is.True(res.ServerTimestamp.Unix() != 0)
	is.True(res.RestServerURL != "")
}

func TestClusterInfo(t *testing.T) {
	is := is.New(t)

	wsTestSrv.SetHandler(func(reqData map[string]interface{}, c *websocket.Conn) map[string]interface{} {
		is.Equal(reqData["action"].(string), "cluster/info")
		return stubs.ResponseStub.ClusterInfo(reqData["requestId"].(string))
	})

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	clusterInfo, dhErr := client.ClusterInfo()

	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
	}

	is.True(clusterInfo != nil)
	is.True(clusterInfo.BootstrapServers != "")
	is.True(clusterInfo.ZookeeperConnect != "")
}
