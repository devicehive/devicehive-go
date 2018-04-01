package dh_ws_test

import (
	"testing"
	"github.com/matryer/is"
)

func TestAuthenticate(t *testing.T) {
	if *tok == "" {
		t.Skip("Access token is not specified, skipping TestAuthenticate")
	}

	is := is.New(t)

	res, err := client.Authenticate(*tok)

	is.True(err == nil)
	is.True(res)
}
