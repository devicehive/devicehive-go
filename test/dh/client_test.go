// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package dh__test

import (
	dh "github.com/devicehive/devicehive-go"
	"github.com/devicehive/devicehive-go/test/stubs"
	"net/http"
	"testing"
)

var testTokens = []byte(`{"accessToken":"test","refreshToken":"test"}`)
var response401 = []byte(`{"timestamp":"2018-05-25T05:20:44.181","status":401,"error":"Unauthorized","message":"Token expired"}`)


func TestReauthorizationByCreds(t *testing.T) {
	httpSrv, httpAddr, httpClose := stubs.StartHTTPTestServer()
	defer httpClose()

	requestCount := 0
	httpSrv.SetRequestHandler(func(reqData map[string]interface{}, rw http.ResponseWriter) {
		requestCount++

		if requestCount == 2 {
			rw.Write(testTokens)
		} else if requestCount == 3 {
			rw.WriteHeader(401)
			rw.Write(response401)
		} else if requestCount == 4 {
			if reqData["login"] == "" || reqData["password"] == "" {
				t.Fatal("Not a token creation by credentials request")
			} else {
				rw.Write(testTokens)
			}
		} else {
			rw.Write([]byte(`{"id":1,"login":"dhadmin"}`))
		}
	})

	client, err := dh.ConnectWithCreds(httpAddr, "dhadmin", "test_password")
	if err != nil {
		t.Fatal(err)
	}

	user, err := client.GetCurrentUser()
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

		if requestCount == 1 {
			rw.WriteHeader(401)
			rw.Write(response401)
		} else if requestCount == 2 {
			if reqData["refreshToken"] == "" {
				t.Fatal("Not a token refresh request")
			} else {
				rw.Write(testTokens)
			}
		} else {
			rw.Write([]byte(`{"id":1,"login":"dhadmin"}`))
		}
	})

	client, err := dh.ConnectWithToken(httpAddr, "test", "test")
	if err != nil {
		t.Fatal(err)
	}

	user, err := client.GetCurrentUser()
	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("User hasn't been obtained")
	}
}
