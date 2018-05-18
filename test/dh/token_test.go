package dh_test

import (
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client := connect(addr)

	userId := 123
	actions := []string{"ManageToken", "ManageNetworks"}
	networkIds := []string{"n1", "n2"}
	deviceTypeIds := []string{"d1", "d2"}
	expiration := time.Now()
	accessToken, refreshToken, dhErr := client.CreateToken(userId, expiration, actions, networkIds, deviceTypeIds)

	if dhErr != nil {
		t.Fatalf("%s: %v", dhErr.Name(), dhErr)
	}

	is.True(accessToken != "")
	is.True(refreshToken != "")
}

func TestRefreshToken(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client := connect(addr)

	accessToken, dhErr := client.RefreshToken()

	if dhErr != nil {
		t.Fatalf("%s: %v", dhErr.Name(), dhErr)
	}

	is.True(accessToken != "")
}
