package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/matryer/is"
	"testing"
)

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
