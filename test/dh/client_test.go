package dh__test

import (
	"testing"
	"github.com/devicehive/devicehive-go/test/stubs"
	"net/http"
	"github.com/devicehive/devicehive-go/dh"
)

func TestReauthorizationByCreds(t *testing.T) {
	httpSrv, httpAddr, httpClose := stubs.StartHTTPTestServer()
	defer httpClose()

	requestCount := 0
	httpSrv.SetRequestHandler(func(reqData map[string]interface{}, rw http.ResponseWriter) {
		requestCount++

		tokens := []byte(`{"accessToken":"test","refreshToken":"test"}`)
		response401 := []byte(`{"timestamp":"2018-05-25T05:20:44.181","status":401,"error":"Unauthorized","message":"Token expired","path":"/api/rest/user"}`)
		if requestCount == 1 {
			rw.Write(tokens)
		} else if requestCount == 2 {
			rw.WriteHeader(401)
			rw.Write(response401)
		} else if requestCount == 3 {
			if reqData["login"] == "" || reqData["password"] == "" {
				t.Fatal("Not a token creation by credentials request")
			} else {
				rw.Write(tokens)
			}
		} else {
			rw.Write([]byte(`{"id":1,"login":"dhadmin"}`))
		}
	})

	newClient, err := dh.ConnectWithCreds(httpAddr, "dhadmin", "test_password")
	if err != nil {
		t.Fatal(err)
	}

	user, err := newClient.GetCurrentUser()
	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("User hasn't been obtained")
	}
}

func TestReauthorizationByRefreshToken(t *testing.T) {
	httpSrv, httpAddr, httpClose := stubs.StartHTTPTestServer()
	defer httpClose()

	requestCount := 0
	httpSrv.SetRequestHandler(func(reqData map[string]interface{}, rw http.ResponseWriter) {
		requestCount++

		tokens := []byte(`{"accessToken":"test","refreshToken":"test"}`)
		response401 := []byte(`{"timestamp":"2018-05-25T05:20:44.181","status":401,"error":"Unauthorized","message":"Token expired","path":"/api/rest/user"}`)
		if requestCount == 1 {
			rw.WriteHeader(401)
			rw.Write(response401)
		} else if requestCount == 2 {
			if reqData["refreshToken"] == "" {
				t.Fatal("Not a token refresh request")
			} else {
				rw.Write(tokens)
			}
		} else {
			rw.Write([]byte(`{"id":1,"login":"dhadmin"}`))
		}
	})

	newClient, err := dh.ConnectWithToken(httpAddr, "test", "test")
	if err != nil {
		t.Fatal(err)
	}

	user, err := newClient.GetCurrentUser()
	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("User hasn't been obtained")
	}
}
