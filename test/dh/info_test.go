package dh_test

import (
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/matryer/is"
	"testing"
)

func TestGetInfo(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client := connect(addr)

	res, err := client.GetInfo()
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	is.True(res != nil)
	is.True(res.APIVersion != "")
	is.True(res.ServerTimestamp.Unix() > 0)
	is.True(res.RestServerURL != "")
}

func TestGetClusterInfo(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client := connect(addr)

	clusterInfo, err := client.GetClusterInfo()
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	is.True(clusterInfo != nil)
	is.True(clusterInfo.BootstrapServers != "")
	is.True(clusterInfo.ZookeeperConnect != "")
}
