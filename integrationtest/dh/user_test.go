package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/matryer/is"
	"testing"
)

func TestUserCreationAndObtaining(t *testing.T) {
	is := is.New(t)

	user, err := client.CreateUser("go-test", "go-test", 1, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = user.Remove()
		if err != nil {
			t.Fatal(err)
		}
	}()

	is.True(user != nil)
	is.True(user.Id != 0)

	sameUser, err := client.GetUser(user.Id)
	if err != nil {
		t.Fatal(err)
	}

	is.True(sameUser != nil)

	currentUser, err := client.GetCurrentUser()
	if err != nil {
		t.Fatal(err)
	}

	is.True(currentUser != nil)
	is.Equal(currentUser.Id, int64(*userId))

	list, err := client.ListUsers(&dh.ListParams{
		UserStatus: 0,
	})
	if err != nil {
		t.Fatal(err)
	}

	is.True(len(list) > 0)
}

func TestUser(t *testing.T) {
	is := is.New(t)

	user, err := client.CreateUser("go-test", "go-test", 1, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = user.Remove()
		if err != nil {
			t.Fatal(err)
		}
	}()

	user.Data = map[string]interface{}{
		"test": "test",
	}

	err = user.Save()
	if err != nil {
		t.Fatal(err)
	}

	err = user.UpdatePassword("brand_new_password")
	if err != nil {
		t.Fatal(err)
	}

	network, err := client.CreateNetwork("go-test-user-network", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		network.Remove()
		if err != nil {
			t.Fatal(err)
		}
	}()

	err = user.AssignNetwork(network.Id)
	if err != nil {
		t.Fatal(err)
	}

	networkList, err := user.ListNetworks()
	if err != nil {
		t.Fatal(err)
	}

	is.Equal(len(networkList), 1)
	is.Equal(networkList[0].Name, "go-test-user-network")
}
