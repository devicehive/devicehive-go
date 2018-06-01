package dh_wsclient_test

import (
	"testing"
	"encoding/json"
	"github.com/devicehive/devicehive-go"
)

func TestWSClientUser(t *testing.T) {
	err := wsclient.CreateUser("go-test", "go-test", 0, nil, false)
	if err != nil {
		t.Fatal(err)
	}

	userId := 0
	testResponse(t, func(data []byte) {
		res := make(map[string]interface{})
		json.Unmarshal(data, &res)
		userId = int(res["id"].(float64))
	})
	defer func() {
		err = wsclient.DeleteUser(userId)
		if err != nil {
			t.Fatal(err)
		}

		testResponse(t, nil)
	}()

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

	err = wsclient.GetCurrentUser()
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)
}

func TestWSClientUserNetworks(t *testing.T) {
	err := wsclient.CreateUser("go-test", "go-test", 0, nil, false)
	if err != nil {
		t.Fatal(err)
	}

	userId := 0
	testResponse(t, func(data []byte) {
		res := make(map[string]interface{})
		json.Unmarshal(data, &res)
		userId = int(res["id"].(float64))
	})
	defer func() {
		err = wsclient.DeleteUser(userId)
		if err != nil {
			t.Fatal(err)
		}

		testResponse(t, nil)
	}()

	err = wsclient.CreateNetwork("test-network", "")
	if err != nil {
		t.Fatal(err)
	}

	networkId := 0
	testResponse(t, func(data []byte) {
		res := make(map[string]interface{})
		json.Unmarshal(data, &res)
		networkId = int(res["id"].(float64))
	})
	defer func() {
		wsclient.DeleteNetwork(networkId)
		testResponse(t, nil)
	}()

	err = wsclient.UserAssignNetwork(userId, networkId)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	err = wsclient.UserUnassignNetwork(userId, networkId)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)
}

func TestWSClientUserDeviceTypes(t *testing.T) {
	err := wsclient.CreateUser("go-test", "go-test", 0, nil, false)
	if err != nil {
		t.Fatal(err)
	}

	userId := 0
	testResponse(t, func(data []byte) {
		res := make(map[string]interface{})
		json.Unmarshal(data, &res)
		userId = int(res["id"].(float64))
	})
	defer func() {
		err = wsclient.DeleteUser(userId)
		if err != nil {
			t.Fatal(err)
		}

		testResponse(t, nil)
	}()

	err = wsclient.CreateDeviceType("test-device-type", "")
	if err != nil {
		t.Fatal(err)
	}

	deviceTypeId := 0
	testResponse(t, func(data []byte) {
		res := make(map[string]interface{})
		json.Unmarshal(data, &res)
		deviceTypeId = int(res["id"].(float64))
	})
	defer func() {
		wsclient.DeleteDeviceType(deviceTypeId)
		testResponse(t, nil)
	}()

	err = wsclient.UserAssignDeviceType(userId, deviceTypeId)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	err = wsclient.ListUserDeviceTypes(userId)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	err = wsclient.UserUnassignDeviceType(userId, deviceTypeId)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	err = wsclient.AllowAllDeviceTypes(userId)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)

	err = wsclient.DisallowAllDeviceTypes(userId)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)
}
