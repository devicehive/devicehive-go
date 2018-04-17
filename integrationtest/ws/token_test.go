package dh_test

import (
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestTokenByPayload(t *testing.T) {
	is := is.New(t)

	expiration := time.Now().Add(1 * time.Second)
	accTok, refTok, err := client.TokenByPayload(1, nil, nil, nil, &expiration)

	is.True(err == nil)
	is.True(accTok != "")
	is.True(refTok != "")
}

func TestRefreshToken(t *testing.T) {
	is := is.New(t)

	accessToken, err := client.RefreshToken()

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}

	is.True(accessToken != "")
}
