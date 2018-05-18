package dh_test

import (
	"testing"
	"github.com/matryer/is"
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

	err = network.Remove()
	if err != nil {
		t.Fatal(err)
	}
}
