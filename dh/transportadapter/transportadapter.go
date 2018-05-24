package transportadapter

import (
	"github.com/devicehive/devicehive-go/internal/transport"
	"time"
)

func New(tsp transport.Transporter) TransportAdapter {
	if tsp.IsWS() {
		return &WSAdapter{
			transport: tsp,
		}
	}

	return &HTTPAdapter{
		transport: tsp,
	}
}

type TransportAdapter interface {
	Request(resourceName string, data map[string]interface{}, timeout time.Duration) (res []byte, err error)
	Subscribe(resourceName string, pollingWaitTimeoutSeconds int, params map[string]interface{}) (tspChan chan []byte, subscriptionId string, err *transport.Error)
	Unsubscribe(resourceName, subscriptionId string, timeout time.Duration) error
	Authenticate(token string, timeout time.Duration) (result bool, err error)
}
