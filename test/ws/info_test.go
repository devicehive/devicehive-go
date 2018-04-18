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

	res, dhErr := client.GetInfo()

	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
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

	clusterInfo, dhErr := client.GetClusterInfo()

	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
	}

	is.True(clusterInfo != nil)
	is.True(clusterInfo.BootstrapServers != "")
	is.True(clusterInfo.ZookeeperConnect != "")
}
