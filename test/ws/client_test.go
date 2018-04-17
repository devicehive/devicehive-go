package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/matryer/is"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	_, addr, srvClose := stubs.StartWSTestServer()
	defer srvClose()

	is := is.New(t)

	client, err := dh.ConnectWithToken(addr, "accTok", "refTok")

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	is.True(client != nil)
}
