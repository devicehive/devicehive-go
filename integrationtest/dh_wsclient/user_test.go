package dh_wsclient_test

import (
	"testing"
	"encoding/json"
	"github.com/devicehive/devicehive-go"
)

func TestWSClientUser(t *testing.T) {
	err := wsclient.CreateUser("go-test", "go-test", 0, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	userId := 0
	testResponse(t, func(data []byte) {
		res := make(map[string]interface{})
		json.Unmarshal(data, &res)
		userId = int(res["id"].(float64))
	})

	user := devicehive_go.User{
		Login: "go-test-updated",
	}
	err = wsclient.UpdateUser(userId, user)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	err = wsclient.GetUser(userId)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	err = wsclient.ListUsers(nil)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	err = wsclient.DeleteUser(userId)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	err = wsclient.GetCurrentUser()
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)
}
