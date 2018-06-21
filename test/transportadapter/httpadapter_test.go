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
