package dh_test

import (
	"testing"
	"github.com/matryer/is"
	"github.com/devicehive/devicehive-go/dh"
)

func TestNetwork(t *testing.T) {
	is := is.New(t)

	network, err := client.CreateNetwork("go-test-network", "go sdk test network")
	if err != nil {
		t.Fatal(err)
	}

	is.True(network != nil)
	is.True(network.Id != 0)

	network.Description = "updated go sdk test network"
	err = network.Save()
	if err != nil {
		t.Fatal(err)
	}

	sameNetwork, err := client.GetNetwork(network.Id)
	if err != nil {
		t.Fatal(err)
	}

	is.True(network != nil)
	is.Equal(sameNetwork.Name, "go-test-network")

	list, err := client.ListNetworks(&dh.ListParams{
		NamePattern: "go-%-network",
	})
	if err != nil {
		t.Fatal(err)
	}

	is.Equal(len(list), 1)
	is.Equal(list[0].Name, "go-test-network")

	err = network.Remove()
	if err != nil {
		t.Fatal(err)
	}
}
