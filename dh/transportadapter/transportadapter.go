package transportadapter

import "github.com/devicehive/devicehive-go/internal/transport"

func New(tsp transport.Transporter) TransportAdapter {
	if tsp.IsWS() {
		return &WSAdapter{}
	}

	return &HTTPAdapter{}
}

type TransportAdapter interface {
	HandleResponseError(rawRes []byte) error
	ResolveResource(resName string, data map[string]interface{}) (resource, method string)
	BuildRequestData(resourceName string, rawData map[string]interface{}) interface{}
}
