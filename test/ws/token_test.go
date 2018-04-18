package dh_test

import (
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestTokenByCreds(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client := connect(addr)

	accessToken, refreshToken, dhErr := client.TokenByCreds("dhadmin", "dhadmin_#911")

	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
	}

	is.True(accessToken != "")
	is.True(refreshToken != "")
}

func TestTokenByPayload(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client := connect(addr)

	userId := 123
	actions := []string{"ManageToken", "ManageNetworks"}
	networkIds := []string{"n1", "n2"}
	deviceTypeIds := []string{"d1", "d2"}
	expiration := time.Now()
	accessToken, refreshToken, dhErr := client.TokenByPayload(userId, actions, networkIds, deviceTypeIds, &expiration)

	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
	}

	is.True(accessToken != "")
	is.True(refreshToken != "")
}

func TestTokenRefresh(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client := connect(addr)

	accessToken, dhErr := client.TokenRefresh("test refresh token")

	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
	}

	is.True(accessToken != "")
}
