package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/matryer/is"
	"testing"
	"github.com/devicehive/devicehive-go/test/stubs"
)

func TestServerInfo(t *testing.T) {
	wsTestSrv := &stubs.WSTestServer{}

	addr := wsTestSrv.Start()
	defer wsTestSrv.Close()

	is := is.New(t)

	client, err := dh.Connect(addr)

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
	wsTestSrv := &stubs.WSTestServer{}

	addr := wsTestSrv.Start()
	defer wsTestSrv.Close()

	is := is.New(t)

	client, err := dh.Connect(addr)

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
