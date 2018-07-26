package transportadapter

import (
	"encoding/json"
	"fmt"
	"github.com/devicehive/devicehive-go"
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/internal/transportadapter"
	"github.com/devicehive/devicehive-go/test/stubs"
	"net/http"
	"testing"
	"time"
)

func TestHTTPSubscriptionLastEntityTimestamp(t *testing.T) {
	httpTestSrv, addr := stubs.StartHTTPTestServer()
	defer httpTestSrv.Close()

	now := &devicehive_go.ISO8601Time{time.Now()}
	commandTimestamp := (&devicehive_go.ISO8601Time{now.Add(1 * time.Second)}).String()
	nextCommandTimestamp := (&devicehive_go.ISO8601Time{now.Add(2 * time.Second)}).String()

	httpTestSrv.SetRequestHandler(func(reqData map[string]interface{}, rw http.ResponseWriter, r *http.Request) {
		timestamp := r.URL.Query()["timestamp"][0]
		if timestamp == now.String() {
			rw.Write([]byte(fmt.Sprintf(`[{"id":1,"timestamp":%q}]`, commandTimestamp)))
		} else if timestamp == commandTimestamp {
			rw.Write([]byte(fmt.Sprintf(`[{"id":2,"timestamp":%q}]`, nextCommandTimestamp)))
		} else {
			rw.Write([]byte(fmt.Sprintf(`[{"id":"success","timestamp":%q}]`, nextCommandTimestamp)))
		}
	})

	httpTsp, err := transport.Create(addr, nil)
	if err != nil {
		t.Fatal(err)
	}

	tspAdapter := transportadapter.New(httpTsp)

	params := map[string]interface{}{
		"timestamp": now.String(),
	}
	subs, _, subsErr := tspAdapter.Subscribe("subscribeCommands", 1, params)
	if subsErr != nil {
		t.Fatal(subsErr)
	}

	lastId := ""
loop:
	for {
		select {
		case d := <-subs.DataChan:
			res := &struct {
				Id json.Number `json:"id"`
			}{}
			json.Unmarshal(d, res)
			if res.Id == "success" {
				break loop
			} else if string(res.Id) == lastId {
				t.Fatal("Timestamp of last entity has not been set for polling")
			}

			lastId = string(res.Id)
		case err := <-subs.ErrChan:
			t.Fatal(err)
		}
	}
}

func TestSubscriptionReauthentication(t *testing.T) {
	httpTestSrv, addr := stubs.StartHTTPTestServer()
	defer httpTestSrv.Close()

	requestsCount := 0
	httpTestSrv.SetRequestHandler(func(reqData map[string]interface{}, rw http.ResponseWriter, r *http.Request) {
		if requestsCount == 0 {
			rw.WriteHeader(401)
			rw.Write([]byte(`{"status": 401,"error": "Unauthorized","message": "Token expired"}`))
		} else if requestsCount == 1 {
			if _, ok := reqData["refreshToken"]; !ok {
				t.Fatal("HTTPAdapter must send token refresh request after it gets 401 Unauthorized")
			}

			rw.WriteHeader(201)
			rw.Write([]byte(`{"accessToken":"test access token"}`))
		} else {
			rw.WriteHeader(200)
			rw.Write([]byte(`[{"id":1,"name":"test command"}]`))
		}

		requestsCount++
	})

	httpTsp, err := transport.Create(addr, nil)
	if err != nil {
		t.Fatal(err)
	}

	tspAdapter := transportadapter.New(httpTsp).(*transportadapter.HTTPAdapter)
	tspAdapter.SetRefreshToken("test refresh token")

	subs, _, tspErr := tspAdapter.Subscribe("subscribeCommands", 1, nil)
	if tspErr != nil {
		t.Fatal(tspErr)
	}

	select {
	case <-subs.DataChan:
	case e := <-subs.ErrChan:
		t.Fatal(e)
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout")
	}
}
