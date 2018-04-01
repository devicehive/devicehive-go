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

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	is.True(res)
}
