package transportadapter

import (
	"github.com/devicehive/devicehive-go/internal/transport"
	"time"
)

func New(tsp transport.Transporter) TransportAdapter {
	adapter := &Adapter{
		transport: tsp,
	}
	if tsp.IsWS() {
		ws := &WSAdapter{
			Adapter: adapter,
		}
		adapter.adapter = ws
		return ws
	}

	http := &HTTPAdapter{
		Adapter: adapter,
	}
	adapter.adapter = http
	return http
}

type TransportAdapter interface {
	Request(resourceName string, data map[string]interface{}, timeout time.Duration) (res []byte, err error)
	Subscribe(resourceName string, pollingWaitTimeoutSeconds int, params map[string]interface{}) (tspChan chan []byte, subscriptionId string, err *transport.Error)
	Unsubscribe(resourceName, subscriptionId string, timeout time.Duration) error
	Authenticate(token string, timeout time.Duration) (result bool, err error)
}
