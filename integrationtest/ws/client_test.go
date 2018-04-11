package dh_test

import (
	"github.com/matryer/is"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	accTok, _, err := client.TokenByCreds(*dhLogin, *dhPass)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	is := is.New(t)

	res, err := client.Authenticate(accTok)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	is.True(res)
}
