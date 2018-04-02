package dh_ws_test

import (
	"github.com/devicehive/devicehive-go/testutils"
	"github.com/matryer/is"
	"testing"
)

func TestServerInfo(t *testing.T) {
	is := is.New(t)

	info, err := client.ServerInfo()

	testutils.LogDHErr(t, err)

	is.True(info != nil)
	is.True(info.APIVersion != "")
	is.True(info.RestServerURL != "")
}

func TestClusterInfo(t *testing.T) {
	is := is.New(t)

	bootstrapServers, zookeeperConnect, err := client.ClusterInfo()

	testutils.LogDHErr(t, err)

	is.True(bootstrapServers != "")
	is.True(zookeeperConnect != "")
}
