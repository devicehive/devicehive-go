package dh_integrationtest

import (
	"flag"
	"github.com/devicehive/devicehive-go/dh"
	"github.com/matryer/is"
	"os"
	"testing"
	"time"
)

const serverAddr = "playground-dev.devicehive.com/api/websocket"
const wsServerAddr = "ws://" + serverAddr

var tok = flag.String("accessToken", "", "Your JWT access token")
var dhLogin = flag.String("dhLogin", "dhadmin", "Your username")
var dhPass = flag.String("dhPassword", "dhadmin_#911", "Your password")

var client *dh.Client

func TestMain(m *testing.M) {
	var err error
	client, err = dh.Connect(wsServerAddr)

	if err != nil {
		panic(err)
	}

	flag.Parse()

	res := m.Run()
	os.Exit(res)
}

func TestAuthenticate(t *testing.T) {
	if *tok == "" {
		t.Skip("Access token is not specified, skipping TestAuthenticate")
	}

	is := is.New(t)

	res, err := client.Authenticate(*tok)

	is.True(err == nil)
	is.True(res)
}

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
