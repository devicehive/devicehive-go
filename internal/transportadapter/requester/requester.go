package requester

import (
	"time"
	"github.com/devicehive/devicehive-go/internal/transport/apirequests"
	"github.com/devicehive/devicehive-go/internal/transport"
)

func New(tsp transport.Transporter) Requester {
	if t, ok := tsp.(*transport.WS); ok {
		return &WSRequester{
			transport: t,
		}
	}

	return nil
}

type Requester interface {
	Request(resourceName string, data map[string]interface{}, timeout time.Duration) ([]byte, error)
	PrepareRequestData(resourceName string, data map[string]interface{}) (resource string, reqParams *apirequests.RequestParams)
}
