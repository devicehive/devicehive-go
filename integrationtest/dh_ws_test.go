package dh_integrationtest

import (
	"testing"
	"github.com/matryer/is"
	"github.com/devicehive/devicehive-go/dh"
	"os"
	"flag"
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

	is.NoErr(err)
	is.True(res)
}

func TestTokenByCreds(t *testing.T) {
	is := is.New(t)

	accTok, refTok, err := client.TokenByCreds(*dhLogin, *dhPass)

	is.NoErr(err)

	is.True(accTok != "")
	is.True(refTok != "")
}

func TestTokenByPayload(t *testing.T) {
	if *tok == "" {
		t.Skip("Access token is not specified, skipping TestTokenByPayload")
	}

	is := is.New(t)

	res, err := client.Authenticate(*tok)

	is.NoErr(err)

	if !res {
		t.Skip("Invalid access token, skipping TestTokenByPayload")
	}

	expiration := time.Now().Add(1 * time.Second)
	accTok, refTok, err := client.TokenByPayload(1, nil, nil, nil, &expiration)

	is.NoErr(err)

	is.True(accTok != "")
	is.True(refTok != "")
}