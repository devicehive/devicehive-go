package transportadapter

import "github.com/devicehive/devicehive-go/internal/transport"

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
}
