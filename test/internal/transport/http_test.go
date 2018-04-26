package transport_test

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
	httpTestSrv, addr, srvClose := stubs.StartHTTPTestServer()
	defer srvClose()

	is := is.New(t)

	httpTestSrv.SetRequestHandler(func(reqData map[string]interface{}, rw http.ResponseWriter) {
		is.True(reqData["requestId"] != "")
		rw.Write([]byte("{}"))
	})

	httpTsp, err := transport.Create(addr)

	is.NoErr(err)

	if err != nil {
		return
	}

	_, tspErr := httpTsp.Request("", map[string]interface{}{}, testHTTPTimeout)

	if tspErr != nil {
		t.Errorf("%s: %v", tspErr.Name(), tspErr)
	}
}

func TestHTTPTimeout(t *testing.T) {
	httpTestSrv, addr, srvClose := stubs.StartHTTPTestServer()
	defer srvClose()

	is := is.New(t)

	httpTestSrv.SetRequestHandler(func(reqData map[string]interface{}, rw http.ResponseWriter) {
		<-time.After(testWSTimeout + 1*time.Second)
		rw.Write([]byte("{\"result\": \"success\"}"))
	})

	httpTsp, err := transport.Create(addr)

	is.NoErr(err)

	res, tspErr := httpTsp.Request("", map[string]interface{}{}, testHTTPTimeout)

	is.True(res == nil)
	is.Equal(tspErr.Name(), transport.TimeoutErr)
}
