// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package transport

import (
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/matryer/is"
	"net/http"
	"testing"
	"time"
)

const testHTTPTimeout = 300 * time.Millisecond

func TestHTTPRequestId(t *testing.T) {
	httpTestSrv, addr := stubs.StartHTTPTestServer()
	defer httpTestSrv.Close()

	is := is.New(t)

	httpTestSrv.SetRequestHandler(func(reqData map[string]interface{}, rw http.ResponseWriter, r *http.Request) {
		is.True(reqData["requestId"] != "")
		rw.Write([]byte("{}"))
	})

	httpTsp, err := transport.Create(addr, nil)

	is.NoErr(err)

	if err != nil {
		return
	}

	_, tspErr := httpTsp.Request("", nil, testHTTPTimeout)

	if tspErr != nil {
		t.Errorf("%s: %v", tspErr.Name(), tspErr)
	}
}

func TestHTTPTimeout(t *testing.T) {
	httpTestSrv, addr := stubs.StartHTTPTestServer()
	defer httpTestSrv.Close()

	is := is.New(t)

	httpTestSrv.SetRequestHandler(func(reqData map[string]interface{}, rw http.ResponseWriter, r *http.Request) {
		<-time.After(testHTTPTimeout + 1*time.Second)
		rw.Write([]byte("{\"result\": \"success\"}"))
	})

	httpTsp, err := transport.Create(addr, nil)

	is.NoErr(err)

	res, tspErr := httpTsp.Request("", nil, testHTTPTimeout)

	is.True(res == nil)
	is.Equal(tspErr.Name(), transport.TimeoutErr)
}

func TestHTTPSubscription(t *testing.T) {
	httpTestSrv, addr := stubs.StartHTTPTestServer()
	defer httpTestSrv.Close()

	is := is.New(t)

	pollReqHandled := false
	httpTestSrv.SetRequestHandler(func(reqData map[string]interface{}, rw http.ResponseWriter, r *http.Request) {
		if pollReqHandled {
			t.Fatal("HTTP transport must stop polling after unsubscription")
		}

		rw.Write([]byte(`[{"id": 1,"command": "command 1"},{"id": 2,"command": "command 2"}]`))
		pollReqHandled = true
	})

	httpTsp, err := transport.Create(addr, nil)
	is.NoErr(err)

	tspChan, subscriptionId, tspErr := httpTsp.Subscribe("device/command/poll?deviceId=device-1", nil)
	if tspErr != nil {
		t.Fatalf("%s: %v", tspErr.Name(), tspErr)
	}

	is.True(subscriptionId != "")

	select {
	case data, ok := <-tspChan.DataChan:
		is.True(ok)
		is.True(data != nil)
	case err := <-tspChan.ErrChan:
		t.Fatal(err)
	case <-time.After(testHTTPTimeout):
		t.Error("subscription event timeout")
	}

	httpTsp.Unsubscribe(subscriptionId)
}
