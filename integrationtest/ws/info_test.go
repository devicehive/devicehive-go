package dh_ws_test

import (
	"testing"
	"github.com/matryer/is"
)

func TestServerInfo(t *testing.T) {
	is := is.New(t)

	info, err := client.ServerInfo()

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	is.True(info != nil)
	is.True(info.APIVersion != "")
	is.True(info.RestServerURL != "")
}

func TestClusterInfo(t *testing.T) {
	is := is.New(t)

	bootstrapServers, zookeeperConnect, err := client.ClusterInfo()

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	is.True(bootstrapServers != "")
	is.True(zookeeperConnect != "")
}
