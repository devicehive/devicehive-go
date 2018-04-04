package dh_test

import (
	"github.com/matryer/is"
	"testing"
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

	info, err := client.ClusterInfo()

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	is.True(info.BootstrapServers != "")
	is.True(info.ZookeeperConnect != "")
}
