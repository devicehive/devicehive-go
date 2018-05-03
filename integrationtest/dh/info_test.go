package dh_test

import (
	"github.com/matryer/is"
	"testing"
)

func TestServerInfo(t *testing.T) {
	is := is.New(t)

	info, err := client.GetInfo()

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	is.True(info != nil)
	is.True(info.APIVersion != "")
	is.True(info.RestServerURL != "" || info.WebSocketServerURL != "")
}

func TestClusterInfo(t *testing.T) {
	is := is.New(t)

	info, err := client.GetClusterInfo()

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	is.True(info.BootstrapServers != "")
	is.True(info.ZookeeperConnect != "")
}
