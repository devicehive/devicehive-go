package dh_ws_test

import (
	"testing"
	"github.com/matryer/is"
	"time"
	"github.com/devicehive/devicehive-go/dh"
)

func TestTokenByCreds(t *testing.T) {
	is := is.New(t)

	accTok, refTok, err := client.TokenByCreds(*dhLogin, *dhPass)

	is.True(err == nil)
	is.True(accTok != "")
	is.True(refTok != "")
}

func TestTokenByPayload(t *testing.T) {
	is := is.New(t)

	accTok, _, err := client.TokenByCreds(*dhLogin, *dhPass)

	res, err := client.Authenticate(accTok)

	is.True(err == nil)

	if !res {
		t.Skip("Invalid access token by credentials, skipping TestTokenByPayload")
	}

	expiration := time.Now().Add(1 * time.Second)
	accTok, refTok, err := client.TokenByPayload(1, nil, nil, nil, &expiration)

	is.True(err == nil)
	is.True(accTok != "")
	is.True(refTok != "")
}

func TestTokenRefresh(t *testing.T) {
	is := is.New(t)

	client, err := dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	_, refTok, err := client.TokenByCreds(*dhLogin, *dhPass)

	accessToken, dhErr := client.TokenRefresh(refTok)

	is.True(dhErr == nil)
	is.True(accessToken != "")
}
