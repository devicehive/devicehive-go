package transport_test

import (
	"time"
	"testing"
	"github.com/devicehive/devicehive-go/test/stubs"
	"github.com/matryer/is"
	"github.com/devicehive/devicehive-go/internal/transport"
	"net/http"
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

	_, tspErr := httpTsp.Request(map[string]interface{}{}, testHTTPTimeout)

	if tspErr != nil {
		t.Errorf("%s: %v", tspErr.Name(), tspErr)
	}
}

