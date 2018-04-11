package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/matryer/is"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	wsTestSrv := &stubs.WSTestServer{}

	addr := wsTestSrv.Start()
	defer wsTestSrv.Close()

	is := is.New(t)

	client, err := dh.Connect(addr)

	if err != nil {
		panic(err)
	}

	res, dhErr := client.Authenticate("someTestToken")

	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
	}

	is.True(res)
}

func TestServiceError(t *testing.T) {
	// @TODO test service error, e.g. 401 Unauthorized
}
