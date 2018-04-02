package dh_ws_test

import (
	"github.com/devicehive/devicehive-go/testutils"
	"github.com/matryer/is"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	if *tok == "" {
		t.Skip("Access token is not specified, skipping TestAuthenticate")
	}

	is := is.New(t)

	res, err := client.Authenticate(*tok)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	is.True(res)
}
