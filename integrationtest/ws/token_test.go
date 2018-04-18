package dh_test

import (
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestTokenByCreds(t *testing.T) {
	is := is.New(t)

	accTok, refTok, err := client.TokenByCreds(*dhLogin, *dhPass)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	is.True(accTok != "")
	is.True(refTok != "")
}

func TestTokenByPayload(t *testing.T) {
	is := is.New(t)

	expiration := time.Now().Add(1 * time.Second)
	accTok, refTok, err := client.TokenByPayload(1, nil, nil, nil, &expiration)

	is.True(err == nil)
	is.True(accTok != "")
	is.True(refTok != "")
}

func TestTokenRefresh(t *testing.T) {
	is := is.New(t)

	_, refTok, err := client.TokenByCreds(*dhLogin, *dhPass)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		t.Skip("Cannot obtain refresh token by credentials, skipping TestTokenRefresh")
	}

	accessToken, err := client.TokenRefresh(refTok)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	is.True(accessToken != "")
}
