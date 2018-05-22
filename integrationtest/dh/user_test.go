package dh_test

import (
	"testing"
	"github.com/matryer/is"
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
}
