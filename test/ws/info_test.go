package dh_ws_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/devicehive/devicehive-go/testutils"
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

	testutils.LogDHErr(t, dhErr)

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

	testutils.LogDHErr(t, dhErr)

	is.True(clusterInfo.BootstrapServers != "")
	is.True(clusterInfo.ZookeeperConnect != "")
}
