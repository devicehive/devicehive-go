package transportadapter

import (
	"github.com/devicehive/devicehive-go/internal/transport"
	"time"
)

func New(tsp transport.Transporter) TransportAdapter {
	if tsp.IsWS() {
		return &WSAdapter{tsp}
	}

	return &HTTPAdapter{tsp}
}

type TransportAdapter interface {
	HandleResponseError(rawRes []byte) error
	ResolveResource(resName string, data map[string]interface{}) (resource, method string)
	BuildRequestData(resourceName string, rawData map[string]interface{}) interface{}
	ExtractResponsePayload(resourceName string, rawRes []byte) []byte
	Request(resourceName, accessToken string, data map[string]interface{}, timeout time.Duration) (res []byte, err error)
	Subscribe(resourceName, accessToken string, pollingWaitTimeoutSeconds int, params map[string]interface{}) (tspChan chan []byte, subscriptionId string, err *transport.Error)
	Unsubscribe(resourceName, accessToken, subscriptionId string, timeout time.Duration) error
}
